package devices

import (
	"log"
	"net/http"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.ListDevices(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, devices)
}

func (h *handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var tempDevice repo.CreateDeviceParams
	if err := json.Read(r, &tempDevice); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdDevice, err := h.service.CreateDevice(r.Context(), tempDevice)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdDevice)
}
