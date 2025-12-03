package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DeviceState string

const (
	DeviceAvailable DeviceState = "available"
	DeviceInUse     DeviceState = "in-use"
	DeviceInactive  DeviceState = "inactive"
)

type Device struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	State     DeviceState
	CreatedAt time.Time `json:"created_at"`
}

func NewDevice(id, name, brand string, state DeviceState, createdAt time.Time) (*Device, error) {

	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	if id == "" {
		id = uuid.New().String()
	}

	device := &Device{
		ID:        id,
		Name:      name,
		Brand:     brand,
		State:     state,
		CreatedAt: createdAt,
	}

	if err := device.Validate(); err != nil {
		return nil, err
	}

	return device, nil
}

func (d *Device) Validate() error {

	_, err := uuid.Parse(d.ID)
	if err != nil {
		return ErrInvalidID
	}

	if d.Name == "" {
		return ErrNameIsRequired
	}
	if d.Brand == "" {
		return ErrBrandIsRequired
	}

	if !d.State.IsValid() {
		if d.State == "" {
			return ErrStateIsRequired
		}
		return ErrInvalidState
	}

	return nil
}

func (d *Device) SetState(s DeviceState) error {
	if !s.IsValid() {
		return ErrInvalidState
	}
	d.State = s
	return nil
}

func (s DeviceState) IsValid() bool {
	switch s {
	case DeviceAvailable, DeviceInUse, DeviceInactive:
		return true
	}
	return false
}

// DeviceRepository defines the interface that the Service layer will use
type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *Device) error
	UpdateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, id string) error
	GetDeviceById(ctx context.Context, id string) (*Device, error)
	GetDevices(ctx context.Context) ([]Device, error)
	GetDevicesByBrand(ctx context.Context, brand string) ([]Device, error)
	GetDevicesByState(ctx context.Context, state string) ([]Device, error)
}
