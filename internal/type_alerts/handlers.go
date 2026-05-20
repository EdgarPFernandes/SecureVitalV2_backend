package type_alerts

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

func (h *handler) ListAlertTypes(w http.ResponseWriter, r *http.Request) {
	typeAlerts, err := h.service.ListTypeAlerts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, typeAlerts)
}

func (h *handler) CreateAlertTypes(w http.ResponseWriter, r *http.Request) {
	var tempTypeAlert repo.TypeAlert
	if err := json.Read(r, &tempTypeAlert); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTypeAlert, err := h.service.CreateAlertType(r.Context(), tempTypeAlert)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdTypeAlert)
}
