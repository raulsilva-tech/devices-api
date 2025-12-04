package domain

import "errors"

var (
	ErrInvalidState      = errors.New("invalid state")
	ErrStateIsRequired   = errors.New("state is required")
	ErrIDIsRequired      = errors.New("id is required")
	ErrBrandIsRequired   = errors.New("brand is required")
	ErrNameIsRequired    = errors.New("name is required")
	ErrInvalidID         = errors.New("invalid uuid")
	ErrDeleteDeviceInUse = errors.New("cannot delete a device in use")
)
