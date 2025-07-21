package auth

import (
	"encoding/json"
	"marketplace/internal/api"
	"marketplace/internal/models"
	"marketplace/internal/validator"
	"marketplace/pkg/logger"

	"net/http"

	"log/slog"
)

func (h *Handler) SingUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = logger.InfoAddValues(ctx,
		logger.HandlerStartedEvent,
		logger.APIMethod, "SingUp",
	)

	var reqData models.AuthRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqData); err != nil {
		api.WriteJSONError(ctx, w, models.ErrorInvalidReqBody)
		return
	}
	r.Body.Close()

	if h.validationOn {
		if err := h.validator.Struct(&reqData); err != nil {
			api.WriteJSONError(ctx, w, validator.ErrorValidate(err))
			return
		}
	}

	userData := models.UserData(reqData)

	token, err := h.authService.SignUp(ctx, userData)
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

	slog.InfoContext(
		ctx,
		logger.HandlerCompletedEvent,
		logger.UserLogin, userData.Login,
		logger.StatusCode, http.StatusOK,
	)
}
