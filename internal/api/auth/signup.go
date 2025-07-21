package auth

import (
	"encoding/json"
	"marketplace/internal/api"
	"marketplace/internal/models"
	"marketplace/internal/validator"
	"marketplace/pkg/logger"

	"net/http"
)

func (h *Handler) SingUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = logger.InfoAddValues(ctx,
		logger.HandlerStartedEvent,
		logger.APIMethod, "SingUp",
	)

	var reqData models.SingUpInRequest

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

	userData := models.UserData(reqData)

	userID, err := h.authService.SignUp(ctx, userData)
	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	token, err := h.authService.GenerateJWT(ctx, userID, userData.Login)
	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)

	logger.InfoAddValues(
		ctx,
		logger.HandlerCompletedEvent,
		logger.UserID, userID,
		logger.StatusCode, http.StatusOK,
	)
}
