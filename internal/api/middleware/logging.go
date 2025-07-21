package middleware

import (
	"marketplace/pkg/logger"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.AddValuesToContext(r.Context(),
			logger.RequestPath, r.URL.Path,
			logger.RequestMethod, r.Method,
			logger.RequestRemoteAddr, r.RemoteAddr,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
