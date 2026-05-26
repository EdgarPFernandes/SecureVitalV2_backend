package patients

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/json"
	"github.com/go-chi/chi/v5"
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

func (h *handler) ListAlertsByPatient(w http.ResponseWriter, r *http.Request) {
	patientIDStr := chi.URLParam(r, "patient_id")

	log.Println("patient_id =", patientIDStr) // debug (optional)

	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid patient id", http.StatusBadRequest)
		return
	}

	alerts, err := h.service.ListAlertsByPatient(r.Context(), patientID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, alerts)
}
