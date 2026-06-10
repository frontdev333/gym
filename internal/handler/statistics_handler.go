package handler

import (
	"net/http"

	"frontdev333/gym/internal/service"

	"github.com/go-chi/chi/v5"
)

type StatisticsHandler struct {
	service *service.StatisticsService
}

func NewStatisticsHandler(svc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{service: svc}
}

func (h *StatisticsHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, err)
		return
	}

	statistics, err := h.service.GetUserStatistics(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, statistics)
}
