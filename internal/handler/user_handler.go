package handler

import (
	"net/http"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *service.UserService
}

type userRequest struct {
	Name string `json:"name"`
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	if users == nil {
		users = []domain.User{}
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	var req userRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.service.Update(r.Context(), userID, req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
