package feed

import (
	"context"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"
	"net/url"
)

func (f *feedService) Feed(ctx context.Context, query url.Values, cursor string, userLogin string) ([]models.AdForFeed, string, error) {
	var (
		ads           []models.AdForFeed
		nextPageURI   string
		newlastAdData models.LastAdData
		err           error
	)

	uriParams, err := f.parseURIParams(query)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Service.Feed.parseURIParams",
			logger.Error, err,
		)
		return ads, nextPageURI, err
	}

	if cursor == "" {
		ads, newlastAdData, err = f.adRepo.Feed(ctx, uriParams, models.LastAdData{}, userLogin)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Service.Feed.UnAuthFeed",
				logger.Error, err,
			)
			return ads, nextPageURI, models.ErrorDb
		}
	}

	if cursor != "" {
		lastAdData, err := parseCursor(f.aesgcm, cursor)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Service.Feed.parseCursor",
				logger.Error, err,
			)
			return ads, nextPageURI, models.ErrorCursorParse
		}
		ads, newlastAdData, err = f.adRepo.Feed(ctx, uriParams, lastAdData, userLogin)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Service.Feed.AuthFeed",
				logger.Error, err,
			)
			return ads, nextPageURI, err
		}
	}

	if newlastAdData == (models.LastAdData{}) {
		return ads, nextPageURI, nil
	}

	newCursor, err := generateCursor(f.aesgcm, newlastAdData)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Service.Feed.generateCursor",
			logger.Error, err,
		)
		return ads, nextPageURI, models.ErrorCursorGenerate
	}

	nextPageURI = generateFeedURI(ctx, f.baseURL, uriParams, newCursor)

	return ads, nextPageURI, nil
}
