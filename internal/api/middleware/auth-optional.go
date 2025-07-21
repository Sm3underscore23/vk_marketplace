package middleware

import (
	"context"
	"log"
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	"marketplace/pkg/logger"
	"net/http"
)

func AuthOptionalMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenCookie, err := r.Cookie("token")
			if err == nil {
				claimData, err := authService.ParseJWT(ctx, "Bearer "+tokenCookie.Value)
				if err == nil {
					ctx = context.WithValue(ctx, models.UserIDKey, claimData.UserId)
					ctx = context.WithValue(ctx, models.UserLoginKey, claimData.UserLogin)

					ctx = logger.AddValuesToContext(
						ctx,
						logger.UserID, claimData.UserId,
						logger.UserLogin, claimData.UserLogin,
					)
				}
				log.Println(err)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
