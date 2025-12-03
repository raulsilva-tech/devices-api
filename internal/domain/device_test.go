package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	//3A
	//arrange
	id, name, brand := uuid.New().String(), "Device 1", "Telec LTDA"
	createdAt := time.Now()
	//act
	d, err := NewDevice(id, name, brand, DeviceAvailable, createdAt)

	//assert
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, d.ID, id)
	assert.Equal(t, d.Name, name)
	assert.Equal(t, d.State, DeviceAvailable)
	assert.Equal(t, d.CreatedAt, createdAt)
}

func TestNewDevice_WhenInvalidId(t *testing.T) {
	//arrange, act
	d, err := NewDevice("sdfs", "Device 1", "Telec LTDA", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrInvalidID)
}

func TestNewDevice_WhenNameIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice(uuid.New().String(), "", "Telec LTDA", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrNameIsRequired)
}

func TestNewDevice_WhenBrandIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice(uuid.New().String(), "Device 1", "", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrBrandIsRequired)
}

func TestNewDevice_WhenStateIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice(uuid.New().String(), "Device 1", "Telec LTDA", "", time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrStateIsRequired)
}

func TestNewDevice_WhenStateIsInvalid(t *testing.T) {
	//arrange, act
	d, err := NewDevice(uuid.New().String(), "Device 1", "Telec LTDA", "invalid", time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrInvalidState)
}
