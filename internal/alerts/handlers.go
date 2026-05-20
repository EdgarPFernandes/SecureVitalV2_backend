package alerts

import (
	"log"
	"net/http"
	"strconv"

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

func (h *handler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.service.ListAlerts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, alerts)
}

func (h *handler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	var tempAlert createAlertParams
	if err := json.Read(r, &tempAlert); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAlert, err := h.service.CreateAlert(r.Context(), tempAlert)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdAlert)
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
