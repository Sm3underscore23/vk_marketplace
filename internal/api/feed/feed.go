package feed

import (
	"encoding/json"
	"log/slog"
	"marketplace/internal/api"
	"marketplace/internal/models"
	"marketplace/pkg/logger"
	"net/http"
)

func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = logger.InfoAddValues(ctx,
		logger.HandlerStartedEvent,
		logger.APIMethod, "Feed",
	)

	query := r.URL.Query()

	cursor := query.Get("cursor")

	userLogin, _ := ctx.Value(models.UserLoginKey).(string)

	ads, nextPageURl, err := h.feedService.Feed(ctx, query, cursor, userLogin)

	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	resp := models.FeedResponse{
		Ads:         ads,
		NextPageURI: nextPageURl,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "can not encode json", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(
		ctx,
		logger.HandlerCompletedEvent,
	)
}
