package advertisement

import (
	"context"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"
)

func (a *adService) CreateAd(ctx context.Context, adData models.AdData) error {
	err := a.adRepo.CreateAd(ctx, adData)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AdvertisementService.CreateAd",
			logger.Error, err,
		)
		return models.ErrorDb
	}

	return err
}
