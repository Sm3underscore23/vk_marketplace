package middleware

import (
	"fmt"
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

func TestAuthStrictMiddleware(t *testing.T) {
	errFormater := func(err error) string { return fmt.Sprintf("{\"errors\":\"%s\"}\n", err) }

	testTable := []struct {
		name           string
		coockieName    string
		coockieValue   string
		mockBehavior   func(c *gomock.Controller) service.AuthService
		expectedStatus int
		expectedData   string
	}{
		{
			name:         "OK",
			coockieName:  "token",
			coockieValue: "test_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().
					ParseJWT(gomock.Any(), "Bearer test_token").
					Return(models.ClaimData{}, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "Invalid cookie name",
			coockieName:  "invalid_name",
			coockieValue: "test_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				return mock_service.NewMockAuthService(c)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedData:   errFormater(models.ErrorInvalidCoockieName),
		},
		{
			name:         "Invalid token",
			coockieName:  "token",
			coockieValue: "invalid_token",
			mockBehavior: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().ParseJWT(gomock.Any(), "Bearer invalid_token").Return(models.ClaimData{}, models.ErrorInvalidJWTFormat)
				return mock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedData:   errFormater(models.ErrorInvalidJWTFormat),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := testCase.mockBehavior(ctrl)
			middleware := AuthStrictMiddleware(authService)

			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, okID := r.Context().Value(models.UserIDKey).(int)
				_, okLogin := r.Context().Value(models.UserLoginKey).(string)
				require.True(t, okID, "userID not found in context")
				require.True(t, okLogin, "userLogin not found in context")
			})

			r := chi.NewRouter()
			r.Use(middleware)
			r.Get("/protected", finalHandler)

			req := httptest.NewRequest("GET", "/protected", nil)
			if testCase.coockieName != "" {
				req.AddCookie(&http.Cookie{
					Name:  testCase.coockieName,
					Value: testCase.coockieValue,
				})
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		})
	}
}
