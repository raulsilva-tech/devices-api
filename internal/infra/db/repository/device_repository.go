package repository

import (
	"context"
	"database/sql"

	"github.com/raulsilva-tech/devices-api/internal/domain"
	"github.com/raulsilva-tech/devices-api/internal/infra/db/sqlc"
)

type DeviceRepository struct {
	db      *sql.DB
	Queries *sqlc.Queries
}

func NewDeviceRepository(dbConn *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		db:      dbConn,
		Queries: sqlc.New(dbConn),
	}
}

func (repo *DeviceRepository) CreateDevice(ctx context.Context, device *domain.Device) error {

	return repo.Queries.CreateDevice(ctx, sqlc.CreateDeviceParams{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     string(device.State),
		CreatedAt: device.CreatedAt,
	})
}

func (repo *DeviceRepository) UpdateDevice(ctx context.Context, device *domain.Device) error {

	return repo.Queries.UpdateDevice(ctx, sqlc.UpdateDeviceParams{
		ID:    device.ID,
		Name:  device.Name,
		Brand: device.Brand,
		State: string(device.State),
	})
}

func (repo *DeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	return repo.Queries.DeleteDevice(ctx, id)
}

func (repo *DeviceRepository) GetDeviceById(ctx context.Context, id string) (*domain.Device, error) {

	devDB, err := repo.Queries.GetDeviceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	device := mapDBToDomainDevice(devDB)
	return &device, nil
}

func (repo *DeviceRepository) GetDevices(ctx context.Context) ([]domain.Device, error) {

	devDBList, err := repo.Queries.GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}

	resultList := make([]domain.Device, len(devDBList))

	for i, devDB := range devDBList {
		resultList[i] = mapDBToDomainDevice(devDB)
	}

	return resultList, nil
}

func (repo *DeviceRepository) GetDevicesByBrand(ctx context.Context, brand string) ([]domain.Device, error) {

	devDBList, err := repo.Queries.GetAllDevicesByBrand(ctx, brand)
	if err != nil {
		return nil, err
	}

	resultList := make([]domain.Device, len(devDBList))

	for i, devDB := range devDBList {
		resultList[i] = mapDBToDomainDevice(devDB)
	}

	return resultList, nil
}

func (repo *DeviceRepository) GetDevicesByState(ctx context.Context, state string) ([]domain.Device, error) {

	devDBList, err := repo.Queries.GetAllDevicesByState(ctx, state)
	if err != nil {
		return nil, err
	}

	resultList := make([]domain.Device, len(devDBList))

	for i, devDB := range devDBList {
		resultList[i] = mapDBToDomainDevice(devDB)
	}

	return resultList, nil
}

func mapDBToDomainDevice(d sqlc.Device) domain.Device {
	return domain.Device{
		ID:        d.ID,
		Name:      d.Name,
		Brand:     d.Brand,
		State:     domain.DeviceState(d.State),
		CreatedAt: d.CreatedAt,
	}
}
