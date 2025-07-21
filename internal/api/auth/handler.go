package auth

import (
	service "marketplace/internal/sevice"
	"marketplace/internal/validator"
)

type Handler struct {
	authService  service.AuthService
	validator    validator.CustomValidator
	validationOn bool
}

func New(authService service.AuthService, validationOn bool, validator validator.CustomValidator) *Handler {
	return &Handler{
		authService:  authService,
		validator:    validator,
		validationOn: validationOn,
	}
}
