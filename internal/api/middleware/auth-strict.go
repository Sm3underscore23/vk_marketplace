package middleware

import (
	"context"
	"marketplace/internal/api"
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	"marketplace/pkg/logger"

	"net/http"
)

func AuthStrictMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				api.WriteJSONError(ctx, w, models.ErrorInvalidCoockieName)
				return
			}

			claimData, err := authService.ParseJWT(ctx, "Bearer "+tokenCookie.Value)
			if err != nil {
				api.WriteJSONError(ctx, w, err)
				return
			}

			ctx = context.WithValue(r.Context(), models.UserIDKey, claimData.UserId)
			ctx = context.WithValue(ctx, models.UserLoginKey, claimData.UserLogin)

			ctx = logger.AddValuesToContext(
				ctx,
				logger.UserID, claimData.UserId,
				logger.UserLogin, claimData.UserLogin,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
