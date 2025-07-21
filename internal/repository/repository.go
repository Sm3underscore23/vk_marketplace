package repository

import (
	"context"
	"marketplace/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login string, passwordHash string) (int, error)
	VerifyUser(ctx context.Context, login string) (models.UserVerify, error)
	IsLoginExists(ctx context.Context, login string) (bool, error)
}

type AdvertisementRepository interface {
	CreateAd(ctx context.Context, adData models.AdData) error
	LastCreatedAd(ctx context.Context) (int, error)
	Feed(
		ctx context.Context, uriParams models.FeedURIParams, lastAdData models.LastAdData, userLogin string,
	) ([]models.AdForFeed, models.LastAdData, error)
}
