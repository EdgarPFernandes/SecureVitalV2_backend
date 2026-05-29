package users

import (
	"log"
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
