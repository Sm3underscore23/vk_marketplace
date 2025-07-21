package advertisement

import (
	"marketplace/internal/repository"
	service "marketplace/internal/sevice"
)

type adService struct {
	adRepo repository.AdvertisementRepository
}

func New(adRepo repository.AdvertisementRepository) service.AdvertisementService {
	return &adService{
		adRepo: adRepo,
	}
}
