package feed

import (
	"crypto/cipher"
	"marketplace/internal/repository"
	service "marketplace/internal/sevice"
)

type feedService struct {
	adRepo       repository.AdvertisementRepository
	aesgcm       cipher.AEAD
	defaultLimit uint64
	baseURL      string
}

func New(adRepo repository.AdvertisementRepository, aesgcm cipher.AEAD, defaultLimit uint64, baseURL string) service.FeedService {
	return &feedService{
		adRepo:       adRepo,
		aesgcm:       aesgcm,
		defaultLimit: defaultLimit,
		baseURL:      baseURL,
	}
}
