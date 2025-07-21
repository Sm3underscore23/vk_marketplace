package auth

import (
	"encoding/json"
	"log/slog"
	"marketplace/internal/api"
	"marketplace/internal/models"
	"marketplace/internal/validator"
	"marketplace/pkg/logger"

	"net/http"
)

func (h *Handler) SingIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = logger.InfoAddValues(ctx,
		logger.HandlerStartedEvent,
		logger.APIMethod, "SingIn",
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

	tokenJWT, err := h.authService.VerifyUser(ctx, userData)
	if err != nil {
		api.WriteJSONError(ctx, w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenJWT,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)

	slog.InfoContext(
		ctx,
		logger.HandlerCompletedEvent,
		logger.UserLogin, userData.Login,
		logger.StatusCode, http.StatusOK,
	)
}
