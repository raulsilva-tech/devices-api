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
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
