package middleware

import (
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	mock_service "marketplace/internal/sevice/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthOptionalMiddleware(t *testing.T) {
	testTable := []struct {
		name               string
		cookieName         string
		cookieValue        string
		mockBehavior       func(c *gomock.Controller) service.AuthService
		expectedAuthStatus bool
	}{
		{
			name:        "Authenticated",
			cookieName:  "token",
			cookieValue: "valid_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().
					ParseJWT(gomock.Any(), "Bearer valid_token").
					Return(models.ClaimData{UserId: 1, UserLogin: "test_login"}, nil)
				return mock
			},
			expectedAuthStatus: true,
		},
		{
			name:        "Invalid cookie name",
			cookieName:  "invalid_cookie",
			cookieValue: "some_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				return mock_service.NewMockAuthService(c)
			},
			expectedAuthStatus: false,
		},
		{
			name:        "Invalid token",
			cookieName:  "token",
			cookieValue: "invalid_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().
					ParseJWT(gomock.Any(), "Bearer invalid_token").
					Return(models.ClaimData{}, models.ErrorInvalidJWTFormat)
				return mock
			},
			expectedAuthStatus: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := testCase.mockBehavior(ctrl)
			middleware := AuthOptionalMiddleware(authService)

			var authStatus bool
			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				_, okID := ctx.Value(models.UserIDKey).(int)
				_, okLogin := ctx.Value(models.UserLoginKey).(string)

				authStatus = okID && okLogin
			})

			r := chi.NewRouter()
			r.Use(middleware)
			r.Get("/optional", finalHandler)

			req := httptest.NewRequest("GET", "/optional", nil)
			if testCase.cookieName != "" {
				req.AddCookie(&http.Cookie{
					Name:  testCase.cookieName,
					Value: testCase.cookieValue,
				})
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, testCase.expectedAuthStatus, authStatus)
		})
	}
}
