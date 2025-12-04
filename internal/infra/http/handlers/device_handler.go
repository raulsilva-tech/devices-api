package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raulsilva-tech/devices-api/internal/domain"
	"github.com/raulsilva-tech/devices-api/internal/dto"
	"github.com/raulsilva-tech/devices-api/internal/service"
)

type DeviceHandler struct {
	Service *service.DeviceService
}

func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		Service: svc,
	}
}

// CreateDevice godoc
// @Summary Create a new device
// @Description Creates a new device and returns its ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body dto.DeviceRequest true "Device payload"
// @Success 201 {object} dto.CreateDeviceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {

	var reqBody dto.DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if reqBody.Name == "" || reqBody.Brand == "" || reqBody.State == "" {
		writeJSONError(w, http.StatusBadRequest, "name, brand and state are required")
		return
	}

	id, err := h.Service.CreateDevice(r.Context(), service.CreateDeviceInput{
		Name:  reqBody.Name,
		Brand: reqBody.Brand,
		State: domain.DeviceState(reqBody.State),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidState) {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("state %s is invalid", reqBody.State))
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateDeviceResponse{
		ID: id,
	})
}

// UpdateDevice godoc
// @Summary Update a device
// @Description Update all fields of a device by ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param request body dto.DeviceRequest true "Update payload"
// @Success 200 {object} dto.UpdateDeviceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /devices/{id} [put]
func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	var reqBody dto.DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if reqBody.Name == "" || reqBody.Brand == "" || reqBody.State == "" {
		writeJSONError(w, http.StatusBadRequest, "id, name, brand and state are required")
		return
	}

	output, err := h.Service.UpdateDevice(r.Context(), service.UpdateDeviceInput{
		ID:    id,
		Name:  reqBody.Name,
		Brand: reqBody.Brand,
		State: domain.DeviceState(reqBody.State),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidState) {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("state %s is invalid", reqBody.State))
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.UpdateDeviceResponse{
		UpdatedFields: output.UpdatedFields,
		IgnoredFields: output.IgnoredFields,
		Device: dto.DeviceResponse{
			ID:        output.Device.ID,
			Name:      output.Device.Name,
			Brand:     output.Device.Brand,
			State:     string(output.Device.State),
			CreatedAt: output.Device.CreatedAt,
		},
	}

	writeJSON(w, http.StatusOK, response)
}

// DeleteDevice godoc
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.Service.DeleteDevice(r.Context(), id)
	if err != nil {
		// If the service returns "not found", send 404 instead of 500
		if errors.Is(err, domain.ErrDeleteDeviceInUse) {
			writeJSONError(w, http.StatusConflict, "device is in use and cannot be deleted")
			return
		}
		if errors.Is(err, service.ErrDeviceNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// GetDeviceByID godoc
// @Summary Get a device by ID
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} dto.DeviceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	device, err := h.Service.GetDeviceById(r.Context(), id)
	if err != nil {

		if errors.Is(err, service.ErrDeviceNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapServiceDeviceToDTO(*device))
}

// GetAllDevices godoc
// @Summary List devices
// @Description Returns all devices, or filter by brand/state query params
// @Tags Devices
// @Produce json
// @Param brand query string false "Filter by brand"
// @Param state query string false "Filter by state"
// @Success 200 {array} dto.DeviceResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /devices [get]
func (h *DeviceHandler) GetAllDevices(w http.ResponseWriter, r *http.Request) {

	brand := r.URL.Query().Get("brand")
	if brand != "" {
		devList, err := h.Service.GetDevicesByBrand(r.Context(), brand)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, processDeviceList(devList))
		return
	}

	state := r.URL.Query().Get("state")
	if state != "" {
		devList, err := h.Service.GetDevicesByState(r.Context(), state)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, processDeviceList(devList))
		return
	}

	devList, err := h.Service.GetDevices(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, processDeviceList(devList))
}

func processDeviceList(devList []service.DeviceOutput) []dto.DeviceResponse {

	resultList := make([]dto.DeviceResponse, len(devList))

	for i, device := range devList {
		resultList[i] = mapServiceDeviceToDTO(device)
	}

	return resultList
}

func mapServiceDeviceToDTO(device service.DeviceOutput) dto.DeviceResponse {
	return dto.DeviceResponse{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     string(device.State),
		CreatedAt: device.CreatedAt,
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
