package advertisement

import (
	"encoding/json"
	"marketplace/internal/api"
	"marketplace/internal/models"
	"marketplace/internal/validator"
	"marketplace/pkg/logger"
	"net/http"
	"time"
)

func (h *Handler) CreateAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = logger.InfoAddValues(ctx,
		logger.HandlerStartedEvent,
		logger.APIMethod, "CreateAd",
	)

	userID, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		api.WriteJSONError(ctx, w, models.ErrorGetUserIDCtx)
		return
	}

	var reqData models.CreateAdRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqData); err != nil {
		api.WriteJSONError(ctx, w, models.ErrorInvalidReqBody)
		return
	}
	r.Body.Close()

	if err := h.validator.Struct(&reqData); err != nil {
		api.WriteJSONError(ctx, w, validator.ErrorValidate(err))
		return
	}

	err := validator.CheckImage(ctx, reqData.ImageUrl)
	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	adData := models.AdData{
		AdInfo: models.AdInfo{
			Title:       reqData.Title,
			Price:       reqData.Price,
			Description: reqData.Description,
			ImageUrl:    reqData.ImageUrl,
		},
		AuthorID:  userID,
		CreatedAt: time.Now(),
	}

	err = h.adService.CreateAd(ctx, adData)
	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	logger.InfoAddValues(
		ctx,
		logger.HandlerCompletedEvent,
		logger.StatusCode, http.StatusOK,
	)
}
