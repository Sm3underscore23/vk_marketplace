package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"marketplace/internal/models"
	"marketplace/pkg/logger"
	"net/http"
)

func WriteJSONError(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := models.ErrorResponse{
		Message: err.Error(),
	}
	statusCode := models.GetStatusCode(err)
	w.WriteHeader(statusCode)

	slog.ErrorContext(ctx, logger.HandlerCompletedEvent,
		logger.StatusCode, statusCode,
		logger.Error, err,
	)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write JSONE: %s", err)
	}
}
