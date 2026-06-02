package users

import (
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"

	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/auth"
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

func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, users)
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload

	if err := json.Read(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.service.RegisterUser(
		r.Context(),
		payload,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, map[string]any{
		"user":  toUserResponse(user),
		"token": token,
	})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

	var payload LoginPayload

	if err := json.Read(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(
		r.Context(),
		payload.Email,
		payload.Password,
	)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *handler) Me(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *handler) Dashboard(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data, err := h.service.GetDashboard(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"elderly": data,
	})
}

func (h *handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	data, err := h.service.GetAdminDashboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, data)
}

func (h *handler) PatientDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {
		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)
		return
	}

	data, err := h.service.GetPatientDashboard(
		r.Context(),
		id,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.Write(
		w,
		http.StatusOK,
		data,
	)
}