package feed

import (
	"crypto/cipher"
	"marketplace/internal/repository"
	service "marketplace/internal/sevice"
)

type feedService struct {
	adRepo  repository.AdvertisementRepository
	aesgcm  cipher.AEAD
	baseURL string
}

func New(adRepo repository.AdvertisementRepository, aesgcm cipher.AEAD, baseURL string) service.FeedService {
	return &feedService{
		adRepo:  adRepo,
		aesgcm:  aesgcm,
		baseURL: baseURL,
	}
}
