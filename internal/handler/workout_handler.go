package handler

import (
	"net/http"
	"time"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/service"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	service *service.WorkoutService
}

type workoutRequest struct {
	ExerciseID  string     `json:"exercise_id"`
	PerformedAt *time.Time `json:"performed_at"`
	Amount      *int64     `json:"amount"`
}

func NewWorkoutHandler(svc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{service: svc}
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	var req workoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	workout, err := h.service.Create(r.Context(), service.CreateWorkoutInput{
		UserID:      userID,
		ExerciseID:  req.ExerciseID,
		PerformedAt: req.PerformedAt,
		Amount:      req.Amount,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, workout)
}

func (h *WorkoutHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	workouts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	if workouts == nil {
		workouts = []domain.Workout{}
	}

	writeJSON(w, http.StatusOK, workouts)
}
