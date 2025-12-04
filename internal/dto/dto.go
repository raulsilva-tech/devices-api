package dto

import "time"

// DeviceRequest represents the payload required to create or update a device
// @Description Device request payload
type DeviceRequest struct {
	Name  string `json:"name" example:"iPhone 13 Pro Max"`
	Brand string `json:"brand" example:"Apple"`
	State string `json:"state" example:"available"`
}

// CreateDeviceResponse represents the response returned after a device is created
// @Description Response containing the created device ID
type CreateDeviceResponse struct {
	ID string `json:"id" example:"49e6d977-58a6-4424-a058-8d025991b325"`
}

// UpdateDeviceResponse represents the result of an update operation
// @Description Summary of updated/ignored fields and the updated device
type UpdateDeviceResponse struct {
	UpdatedFields []string       `json:"updated_fields" example:"[\"name\", \"brand\"]"`
	IgnoredFields []string       `json:"ignored_fields" example:"[\"state\"]"`
	Device        DeviceResponse `json:"device"`
}

// DeviceResponse represents a device stored in the system
// @Description Device full information
type DeviceResponse struct {
	ID        string    `json:"id" example:"49e6d977-58a6-4424-a058-8d025991b325"`
	Name      string    `json:"name" example:"Galaxy S21"`
	Brand     string    `json:"brand" example:"Samsung"`
	State     string    `json:"state" example:"in-use"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-10T15:04:05Z"`
}

// ErrorResponse represents an error message
// @Description Error response container
type ErrorResponse struct {
	Error string `json:"error" example:"error description"`
}
