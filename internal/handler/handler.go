package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"frontdev333/gym/internal/domain"

	"github.com/google/uuid"
)

const maxBodySize = 1 << 20

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("writeJSON", "error", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	status, message := mapError(err)
	writeJSON(w, status, errorResponse{
		Error: errorBody{
			Message: message,
			Code:    status,
		},
	})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "resource not found"
	case errors.Is(err, domain.ErrValidation):
		var validationErr domain.ValidationError
		if errors.As(err, &validationErr) {
			return http.StatusBadRequest, validationErr.Message
		}

		return http.StatusBadRequest, "validation error"
	default:
		slog.Error("internal error", "error", err)
		return http.StatusInternalServerError, "internal server error"
	}
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxBodySize))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return domain.NewValidationError("request body is required")
		}

		return domain.NewValidationError("invalid JSON body")
	}

	return nil
}

func parseUUID(value string) (string, error) {
	if value == "" {
		return "", domain.NewValidationError("id is required")
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return "", domain.NewValidationError("invalid id format")
	}

	return parsed.String(), nil
}

func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
