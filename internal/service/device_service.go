package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raulsilva-tech/devices-api/internal/domain"
)

type DeviceService struct {
	repo domain.DeviceRepository
}

func NewDeviceService(repo domain.DeviceRepository) *DeviceService {
	return &DeviceService{
		repo: repo,
	}
}

type CreateDeviceInput struct {
	Name  string
	Brand string
	State domain.DeviceState
}

type UpdateDeviceInput struct {
	ID    string
	Name  string
	Brand string
	State domain.DeviceState
}

type UpdateDeviceOutput struct {
	UpdatedFields []string
	IgnoredFields []string
	Device        DeviceOutput
}

type DeviceOutput struct {
	ID        string
	Name      string
	Brand     string
	State     domain.DeviceState
	CreatedAt time.Time
}

func (s *DeviceService) CreateDevice(ctx context.Context, input CreateDeviceInput) (string, error) {

	device, err := domain.NewDevice(uuid.New().String(), input.Name, input.Brand, input.State, time.Now())
	if err != nil {
		return "", err
	}
	id, err := s.repo.CreateDevice(ctx, device)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, input UpdateDeviceInput) (*UpdateDeviceOutput, error) {

	// • Creation time cannot be updated: UpdateDeviceInput does not offer createdAt field be changed

	// getting device by id to check state
	device, err := s.repo.GetDeviceById(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if device.State != input.State && !input.State.IsValid() {
		return nil, domain.ErrInvalidState
	}

	output := &UpdateDeviceOutput{
		UpdatedFields: []string{},
		IgnoredFields: []string{},
	}

	// • Name and brand properties cannot be updated if the device is in use.
	if device.State == domain.DeviceInUse {

		// Only Name can change
		if input.Name != device.Name {
			device.Name = input.Name
			output.UpdatedFields = append(output.UpdatedFields, "name")
		}

		if input.Brand != device.Brand {
			output.IgnoredFields = append(output.IgnoredFields, "brand")
		}

		if input.State != device.State {
			output.IgnoredFields = append(output.IgnoredFields, "state")
		}

	} else {
		if input.Name != device.Name {
			device.Name = input.Name
			output.UpdatedFields = append(output.UpdatedFields, "name")
		}

		if input.Brand != device.Brand {
			device.Brand = input.Brand
			output.UpdatedFields = append(output.UpdatedFields, "brand")
		}

		if input.State != device.State {
			device.State = input.State
			output.UpdatedFields = append(output.UpdatedFields, "state")
		}
	}

	err = s.repo.UpdateDevice(ctx, device)
	if err != nil {
		return nil, err
	}

	output.Device = DeviceOutput{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     device.State,
		CreatedAt: device.CreatedAt,
	}

	return output, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {

	// getting device by id to check state
	device, err := s.repo.GetDeviceById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("device id %s not found", id)
		}
		return err
	}

	if device.State == domain.DeviceInUse {
		return domain.ErrDeleteDeviceInUse

	}

	return s.repo.DeleteDevice(ctx, id)
}

func (s *DeviceService) GetDeviceById(ctx context.Context, id string) (*DeviceOutput, error) {

	device, err := s.repo.GetDeviceById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device id %s not found", id)
		}
		return nil, err
	}

	return &DeviceOutput{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     device.State,
		CreatedAt: device.CreatedAt,
	}, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]DeviceOutput, error) {
	devList, err := s.repo.GetDevices(ctx)
	if err != nil {
		return []DeviceOutput{}, err
	}
	return processDeviceList(devList)
}

func (s *DeviceService) GetDevicesByBrand(ctx context.Context, brand string) ([]DeviceOutput, error) {
	devList, err := s.repo.GetDevicesByBrand(ctx, brand)
	if err != nil {
		return []DeviceOutput{}, err
	}
	return processDeviceList(devList)
}

func (s *DeviceService) GetDevicesByState(ctx context.Context, state string) ([]DeviceOutput, error) {
	devList, err := s.repo.GetDevicesByState(ctx, state)
	if err != nil {
		return []DeviceOutput{}, err
	}
	return processDeviceList(devList)
}

func processDeviceList(devList []domain.Device) ([]DeviceOutput, error) {

	if len(devList) == 0 {
		return []DeviceOutput{}, nil
	}

	resultList := make([]DeviceOutput, len(devList))

	for i, dev := range devList {
		resultList[i] = mapDomainToServiceDevice(dev)
	}

	return resultList, nil
}

func mapDomainToServiceDevice(device domain.Device) DeviceOutput {
	return DeviceOutput{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     device.State,
		CreatedAt: device.CreatedAt,
	}
}
