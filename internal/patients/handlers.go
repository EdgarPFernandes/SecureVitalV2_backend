package patients

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

func (h *handler) ListPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.service.ListPatients(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, patients)
}

func (h *handler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var tempPatientParams repo.CreatePatientParams
	if err := json.Read(r, &tempPatientParams); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdPatient, err := h.service.CreatePatient(r.Context(), tempPatientParams)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdPatient)
}
