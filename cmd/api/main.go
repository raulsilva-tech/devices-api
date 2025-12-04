package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/raulsilva-tech/devices-api/internal/docs"
	"github.com/raulsilva-tech/devices-api/internal/infra/db/repository"
	"github.com/raulsilva-tech/devices-api/internal/infra/http/handlers"
	"github.com/raulsilva-tech/devices-api/internal/infra/http/middleware"
	"github.com/raulsilva-tech/devices-api/internal/service"
	"github.com/raulsilva-tech/devices-api/shared/env"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	WebServerPort  = env.GetInt("WEBSERVER_PORT", 8080)
	DBPort         = env.GetInt("DB_PORT", 5432)
	DBDriver       = env.GetString("DB_DRIVER", "postgres")
	DBUser         = env.GetString("DB_USER", "myuser")
	DBPassword     = env.GetString("DB_PASSWORD", "mypassword")
	DBHost         = env.GetString("DB_HOST", "localhost")
	DBDatabaseName = env.GetString("DB_NAME", "devices-api")
)

// @title Devices API
// @version 1.0
// @description API for managing devices
// @host localhost:8080
// @BasePath /
func main() {

	dbAddr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBDatabaseName)

	db, err := sql.Open(DBDriver, dbAddr)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	repo := repository.NewDeviceRepository(db)
	svc := service.NewDeviceService(repo)
	devHandler := handlers.NewDeviceHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /devices", devHandler.CreateDevice)
	mux.HandleFunc("PUT /devices/{id}", devHandler.UpdateDevice)
	mux.HandleFunc("DELETE /devices/{id}", devHandler.DeleteDevice)
	mux.HandleFunc("GET /devices/{id}", devHandler.GetDeviceByID)
	mux.HandleFunc("GET /devices", devHandler.GetAllDevices)

	// swagger ui
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	var handler http.Handler = mux
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Timeout(10 * time.Second)(handler)
	handler = middleware.Recover(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", WebServerPort),
		Handler: handler,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Starting API WebServer")
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {

	case err := <-serverErrors:
		log.Println("server error: ", err.Error())

	case sig := <-shutdown:

		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Println("could not shutdown gracefully", err.Error())
			server.Close()
		}
		db.Close()
	}
}
