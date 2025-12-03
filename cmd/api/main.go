package main

import (
	"database/sql"
	"fmt"

	"github.com/raulsilva-tech/devices-api/shared/env"
)

var (
	WebServerPort  = env.GetString("WEBSERVER_PORT", ":8081")
	DBPort         = env.GetString("DB_Port", "5432")
	DBDriver       = env.GetString("DB_DRIVER", "postgres")
	DBUser         = env.GetString("DB_USER", "myuser")
	DBPassword     = env.GetString("DB_PASSWORD", "mypassword")
	DBHost         = env.GetString("DB_HOST", "localhost")
	DBDatabaseName = env.GetString("DB_NAME", "stockdb")
)

func main() {

	dbAddr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBDatabaseName)

	db, err := sql.Open(DBDriver, dbAddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	

}
