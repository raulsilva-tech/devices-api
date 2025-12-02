package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	//3A
	//arrange
	id := "1sdfa"
	name := "Device 1"
	brand := "Telec LTDA"
	createdAt := time.Now()
	//act
	d, err := NewDevice(id, name, brand, DeviceAvailable, createdAt)

	//assert
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, d.ID, id)
	assert.Equal(t, d.Name, name)
	assert.Equal(t, d.state, DeviceAvailable)
	assert.Equal(t, d.CreatedAt, createdAt)
}

func TestNewDevice_WhenIdIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice("", "Device 1", "Telec LTDA", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrIDIsRequired)
}

func TestNewDevice_WhenNameIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice("sdfsd", "", "Telec LTDA", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrNameIsRequired)
}

func TestNewDevice_WhenBrandIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice("sdfs", "Device 1", "", DeviceAvailable, time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrBrandIsRequired)
}

func TestNewDevice_WhenStateIsRequired(t *testing.T) {
	//arrange, act
	d, err := NewDevice("sdfs", "Device 1", "Telec LTDA", "", time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrStateIsRequired)
}

func TestNewDevice_WhenStateIsInvalid(t *testing.T) {
	//arrange, act
	d, err := NewDevice("sdfs", "Device 1", "Telec LTDA", "invalid", time.Now())

	//assert
	assert.NotNil(t, err)
	assert.Nil(t, d)
	assert.Equal(t, err, ErrInvalidState)
}
