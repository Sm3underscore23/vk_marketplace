package feed

// func (f *feedService) UnAuthUserFeed(ctx context.Context, cursor string, uriParams models.FeedURIParams) ([]models.UnAuthUserAdForFeed, string, error) {
// 	var (
// 		feedData  models.UnAuthUserFeed
// 		newCursor string
// 	)
// 	adID, err := parseCursor(f.aesgcm, cursor)
// 	if err != nil {
// 		return feedData.Ads, newCursor, models.ErrorCursorParse
// 	}

// 	feedData, err = f.adRepo.UnAuthUserFeed(ctx, adID, uriParams)
// 	if err != nil {
// 		return feedData.Ads, newCursor, models.ErrorDb
// 	}

// 	if feedData.LastAdID == 0 {
// 		return feedData.Ads, newCursor, nil
// 	}

// 	newCursor, err = generateCursor(f.aesgcm, feedData.LastAdID)
// 	if err != nil {
// 		return feedData.Ads, newCursor, models.ErrorCursorGenerate
// 	}

// 	return feedData.Ads, newCursor, nil
// }
