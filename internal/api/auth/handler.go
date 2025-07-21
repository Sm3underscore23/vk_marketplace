package auth

import (
	service "marketplace/internal/sevice"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	authService service.AuthService
	validator   *validator.Validate
}

func New(authService service.AuthService, validator *validator.Validate) *Handler {
	return &Handler{
		authService: authService,
		validator:   validator,
	}
}
