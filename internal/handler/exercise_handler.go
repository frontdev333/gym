package handler

import (
	"net/http"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/service"

	"github.com/go-chi/chi/v5"
)

type ExerciseHandler struct {
	service *service.ExerciseService
}

type exerciseRequest struct {
	Title string `json:"title"`
}

func NewExerciseHandler(svc *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{service: svc}
}

func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req exerciseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	exercise, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, exercise)
}

func (h *ExerciseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	if exercises == nil {
		exercises = []domain.Exercise{}
	}

	writeJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := parseUUID(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	exercise, err := h.service.GetByID(r.Context(), exerciseID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, exercise)
}

func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := parseUUID(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	var req exerciseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	exercise, err := h.service.Update(r.Context(), exerciseID, req.Title)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, exercise)
}

func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := parseUUID(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), exerciseID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
