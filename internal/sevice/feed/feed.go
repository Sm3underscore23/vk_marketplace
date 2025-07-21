package feed

import (
	"context"
	"marketplace/internal/models"
)

func (f *feedService) Feed(ctx context.Context, uriParams models.FeedURIParams, cursor string, userLogin string) ([]models.AdForFeed, string, error) {
	var (
		ads           []models.AdForFeed
		nextPageURI   string
		newlastAdData models.LastAdData
		err           error
	)

	if cursor == "" {
		ads, newlastAdData, err = f.adRepo.Feed(ctx, uriParams, models.LastAdData{}, userLogin)
		if err != nil {
			return ads, nextPageURI, err
		}
	}

	if cursor != "" {
		lastAdData, err := parseCursor(f.aesgcm, cursor)
		if err != nil {
			return ads, nextPageURI, models.ErrorCursorParse
		}
		ads, newlastAdData, err = f.adRepo.Feed(ctx, uriParams, lastAdData, userLogin)
		if err != nil {
			return ads, nextPageURI, err
		}
	}

	if newlastAdData == (models.LastAdData{}) {
		return ads, nextPageURI, nil
	}

	newCursor, err := generateCursor(f.aesgcm, newlastAdData)
	if err != nil {
		return ads, nextPageURI, models.ErrorCursorGenerate
	}

	nextPageURI = generateFeedURI(ctx, f.baseURL, uriParams, newCursor)

	return ads, nextPageURI, nil
}
