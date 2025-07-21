package advertisement

import (
	service "marketplace/internal/sevice"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	adService service.AdvertisementService
	validator *validator.Validate
}

func New(adService service.AdvertisementService, validator *validator.Validate) *Handler {
	return &Handler{
		adService: adService,
		validator: validator,
	}
}
