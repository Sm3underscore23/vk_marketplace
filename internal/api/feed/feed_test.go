package feed

import (
	"context"
	"fmt"
	"marketplace/internal/models"
	service "marketplace/internal/sevice"
	mock_service "marketplace/internal/sevice/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFeed(t *testing.T) {
	errFormater := func(err error) string { return fmt.Sprintf("{\"errors\":\"%s\"}\n", err) }

	testTable := []struct {
		name            string
		mockFeedService func(c *gomock.Controller, query url.Values) service.FeedService
		expectedStatus  int
		expectedData    string
	}{
		{
			name: "OK",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "", nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedData:   "{\"ads\":null,\"next_page_uri\":\"\"}\n",
		},
		{
			name: "failed to parse uri: unknown param",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorUnknownURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorUnknownURIParam),
		},
		{
			name: "failed to parse uri: unknown param",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidSordByURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidSordByURIParam),
		},
		{
			name: "failed to parse uri: unknown param",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorUnknownURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorUnknownURIParam),
		},
		{
			name: "failed to parse uri: invalid sort_by",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidSordByURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidSordByURIParam),
		},
		{
			name: "failed to parse uri: invalid order",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidOrderURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidOrderURIParam),
		},
		{
			name: "failed to parse uri: invalid price_min",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidPriceMinURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidPriceMinURIParam),
		},
		{
			name: "failed to parse uri: invalid price_max",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidPriceMaxURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   errFormater(models.ErrorInvalidPriceMaxURIParam),
		},
		{
			name: "failed to parse uri: price_min > price_max",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorInvalidPricesURIParam,
				)
				return mock
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   "{\"errors\":\"price filter failed: price_min \\u003e price_max\"}\n",
		},
		{
			name: "error db",
			mockFeedService: func(c *gomock.Controller, query url.Values) service.FeedService {
				mock := mock_service.NewMockFeedService(c)
				mock.EXPECT().Feed(gomock.Any(), query, "", "").Return(nil, "",
					models.ErrorDb,
				)
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

			req := httptest.NewRequest("POST", "/feed", nil)

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			handler := New(testCase.mockFeedService(ctrl, req.URL.Query()))

			handler.Feed(w, req)

			require.Equal(t, testCase.expectedData, w.Body.String())
			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		})
	}
}
