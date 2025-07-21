package service

import (
	"context"
	"marketplace/internal/models"
	"net/url"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type AuthService interface {
	VerifyUser(ctx context.Context, userData models.UserData) (string, error)
	SignUp(ctx context.Context, userData models.UserData) (string, error)
	ParseJWT(ctx context.Context, authHeader string) (models.ClaimData, error)
}

type AdvertisementService interface {
	CreateAd(ctx context.Context, adData models.AdData) error
}

type FeedService interface {
	Feed(ctx context.Context, query url.Values, cursor string, userLogin string) ([]models.AdForFeed, string, error)
}
