package advertisement

import (
	service "marketplace/internal/sevice"
	"marketplace/internal/validator"
)

type Handler struct {
	adService    service.AdvertisementService
	validator    validator.CustomValidator
	validationOn bool
}

func New(adService service.AdvertisementService, validationOn bool, validator validator.CustomValidator) *Handler {
	return &Handler{
		adService:    adService,
		validationOn: validationOn,
		validator:    validator,
	}
}
