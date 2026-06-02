package dashboard

import (
	"net/http"

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

func (h *handler) AdminDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	data, err := h.service.GetAdminDashboard(
		r.Context(),
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.Write(w, http.StatusOK, data)
}