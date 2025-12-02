package domain

import "time"

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
	state     DeviceState
	CreatedAt time.Time `json:"created_at"`
}

func NewDevice(id, name, brand string, state DeviceState, createdAt time.Time) (*Device, error) {

	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	device := &Device{
		ID:        id,
		Name:      name,
		Brand:     brand,
		state:     state,
		CreatedAt: createdAt,
	}

	if err := device.Validate(); err != nil {
		return nil, err
	}

	return device, nil
}

func (d *Device) Validate() error {

	if d.ID == "" {
		return ErrIDIsRequired
	}
	if d.Name == "" {
		return ErrNameIsRequired
	}
	if d.Brand == "" {
		return ErrBrandIsRequired
	}

	if !d.state.IsValid() {
		if d.state == "" {
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
	d.state = s
	return nil
}

func (s DeviceState) IsValid() bool {
	switch s {
	case DeviceAvailable, DeviceInUse, DeviceInactive:
		return true
	}
	return false
}
