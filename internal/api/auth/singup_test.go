package auth

import (
	"context"
	"fmt"
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	mock_service "marketplace/internal/sevice/mocks"
	"marketplace/internal/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSingUp(t *testing.T) {
	errFormater := func(err error) string { return fmt.Sprintf("{\"errors\":\"%s\"}\n", err) }

	reqData := models.UserData{
		Login:    "test_login",
		Password: "tets_password",
	}
	validReqData := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, reqData.Login, reqData.Password)

	validator := validator.New(0)

	testTable := []struct {
		name            string
		reqData         string
		mockAuthService func(c *gomock.Controller) service.AuthService
		expectedStatus  int
		expectedData    string
	}{
		{
			name:    "OK",
			reqData: validReqData,
			mockAuthService: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().SignUp(gomock.Any(), reqData).Return("", nil)
				return mock
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "bad request: invalid data",
			reqData: "invalid_data",
			mockAuthService: func(c *gomock.Controller) service.AuthService {
				return mock_service.NewMockAuthService(c)
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidReqBody),
		},
		{
			name:    "bad request: login already exists",
			reqData: validReqData,
			mockAuthService: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().SignUp(gomock.Any(), reqData).Return("", models.ErrorLoginAlreadyExists)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorLoginAlreadyExists),
		},
		{
			name:    "error db",
			reqData: validReqData,
			mockAuthService: func(c *gomock.Controller) service.AuthService {
				mock := mock_service.NewMockAuthService(c)
				mock.EXPECT().SignUp(gomock.Any(), reqData).Return("", models.ErrorDb)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   errFormater(models.ErrorDb),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			req := httptest.NewRequest("POST", "/sign_in", strings.NewReader(testCase.reqData))

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			handler := New(testCase.mockAuthService(ctrl), false, validator)

			handler.SingUp(w, req)

			require.Equal(t, testCase.expectedData, w.Body.String())
			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		})
	}
}
