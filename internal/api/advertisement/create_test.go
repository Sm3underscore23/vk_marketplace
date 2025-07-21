package advertisement

import (
	"context"
	"fmt"
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	mock_service "marketplace/internal/sevice/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	monkey "bou.ke/monkey"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAd(t *testing.T) {

	monkey.Patch(time.Now, func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) })

	errFormater := func(err error) string { return fmt.Sprintf("{\"errors\":\"%s\"}\n", err) }

	validAdInfo := models.AdInfo{
		Title:       "test title",
		Description: "test description",
		Price:       1234,
		ImageUrl:    "https://raw.githubusercontent.com/mohammadimtiazz/standard-test-images-for-Image-Processing/refs/heads/master/standard_test_images/tulips.png",
	}

	validReqData := fmt.Sprintf("{\"title\":\"%s\",\"description\":\"%s\",\"price\":%d,\"image_url\":\"%s\"}",
		validAdInfo.Title,
		validAdInfo.Description,
		validAdInfo.Price,
		validAdInfo.ImageUrl,
	)

	testTable := []struct {
		name           string
		withoutUserId  bool
		reqData        string
		mockAdService  func(c *gomock.Controller, userID int) service.AdvertisementService
		expectedStatus int
		expectedData   string
	}{
		{
			name:    "OK",
			reqData: validReqData,
			mockAdService: func(c *gomock.Controller, userID int) service.AdvertisementService {
				mock := mock_service.NewMockAdvertisementService(c)
				mock.EXPECT().CreateAd(gomock.Any(), models.AdData{
					AdInfo:    validAdInfo,
					AuthorID:  1,
					CreatedAt: time.Now(),
				}).Return(nil)
				return mock
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "bad request: invalid body",
			reqData: "invalid_data",
			mockAdService: func(c *gomock.Controller, userID int) service.AdvertisementService {
				return mock_service.NewMockAdvertisementService(c)
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidReqBody),
		},
		{
			name:          "user_id not found",
			withoutUserId: true,
			reqData:       "invalid_data",
			mockAdService: func(c *gomock.Controller, userID int) service.AdvertisementService {
				return mock_service.NewMockAdvertisementService(c)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   errFormater(models.ErrorGetUserIDCtx),
		},
		{
			name:    "db error",
			reqData: validReqData,
			mockAdService: func(c *gomock.Controller, userID int) service.AdvertisementService {
				mock := mock_service.NewMockAdvertisementService(c)
				mock.EXPECT().CreateAd(gomock.Any(), models.AdData{
					AdInfo:    validAdInfo,
					AuthorID:  1,
					CreatedAt: time.Now(),
				}).Return(models.ErrorDb)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   errFormater(models.ErrorDb),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			req := httptest.NewRequest("POST", "/advertisements/create", strings.NewReader(testCase.reqData))

			if !testCase.withoutUserId {
				ctx = context.WithValue(req.Context(), models.UserIDKey, 1)
			}

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := New(testCase.mockAdService(controller, 1), false, nil)

			handler.CreateAd(w, req)

			require.Equal(t, testCase.expectedData, w.Body.String())
			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		})
	}
}
