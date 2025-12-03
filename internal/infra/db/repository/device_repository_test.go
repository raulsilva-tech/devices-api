package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/raulsilva-tech/devices-api/internal/domain"
	"github.com/stretchr/testify/suite"
)

type DeviceRepositoryTestSuite struct {
	DB  *sql.DB
	ctx context.Context
	suite.Suite
}

func TestDeviceRepositorySuite(t *testing.T) {
	suite.Run(t, new(DeviceRepositoryTestSuite))
}

func (suite *DeviceRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *DeviceRepositoryTestSuite) SetupSuite() {
	dbConn, err := migrateDB()
	suite.NoError(err)
	suite.DB = dbConn
	suite.ctx = context.Background()
}

func (suite *DeviceRepositoryTestSuite) TestCreate() {

	_, err := domain.NewDevice(uuid.New().String(), "Device", "Brand", domain.DeviceAvailable, time.Now())
	suite.NoError(err)

}

func (suite *DeviceRepositoryTestSuite) TestUpdate() {

	repo, d, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	d.Name = "Updated Device"
	err = repo.UpdateDevice(suite.ctx, d)
	suite.NoError(err)

	dbDevice, err := repo.GetDeviceById(suite.ctx, d.ID)
	suite.NoError(err)
	suite.Equal(dbDevice.Name, d.Name)

}

func (suite *DeviceRepositoryTestSuite) TestGetByID() {

	repo, d, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	suite.NoError(err)
	dbDevice, err := repo.GetDeviceById(suite.ctx, d.ID)
	suite.NoError(err)
	suite.Equal(dbDevice.Name, d.Name)

}

func (suite *DeviceRepositoryTestSuite) TestDelete() {

	repo, d, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	suite.NoError(err)
	err = repo.DeleteDevice(suite.ctx, d.ID)
	suite.NoError(err)

}

func (suite *DeviceRepositoryTestSuite) TestGetAll() {

	repo, _, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	deviceList, err := repo.Queries.GetAllDevices(suite.ctx)
	suite.NoError(err)
	suite.NotEmpty(deviceList)
	suite.Equal(len(deviceList), 1)
}

func (suite *DeviceRepositoryTestSuite) TestGetAllByBrand() {

	repo, _, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	deviceList, err := repo.Queries.GetAllDevices(suite.ctx)
	suite.NoError(err)
	suite.NotEmpty(deviceList)
	suite.Equal(len(deviceList), 1)
}

func (suite *DeviceRepositoryTestSuite) TestGetAllByState() {

	repo, _, err := createDevice(suite.ctx, suite.DB)
	suite.NoError(err)

	deviceList, err := repo.Queries.GetAllDevices(suite.ctx)
	suite.NoError(err)
	suite.NotEmpty(deviceList)
	suite.Equal(len(deviceList), 1)
}

func createDevice(ctx context.Context, db *sql.DB) (*DeviceRepository, *domain.Device, error) {
	device, err := domain.NewDevice(uuid.New().String(), "Device", "Brand", domain.DeviceAvailable, time.Now())
	if err != nil {
		return nil, nil, err
	}

	repo := NewDeviceRepository(db)
	err = repo.CreateDevice(ctx, device)
	if err != nil {
		return nil, nil, err
	}

	return repo, device, err
}

func migrateDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE devices (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    brand      TEXT NOT NULL,
    state      TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`)

	return db, err
}
