package dto

import "time"

type DeviceRequest struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	State string `json:"state"`
}

type CreateDeviceResponse struct {
	ID string `json:"id"`
}

type UpdateDeviceResponse struct {
	UpdatedFields []string       `json:"updated_fields"`
	IgnoredFields []string       `json:"ignored_fields"`
	Device        DeviceResponse `json:"device"`
}

type DeviceResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
