package validator

import (
	"marketplace/internal/models"

	"github.com/go-playground/validator/v10"
)

var errValidate = map[string]error{
	"login":       models.ErrorLoginValidation,
	"password":    models.ErrorPasswordValidation,
	"title":       models.ErrorTitleValidation,
	"description": models.ErrorDescriptionValidation,
	"price":       models.ErrorPriceValidation,
	"image_url": models.ErrorPriceValidation,
}

func ErrorValidate(err error) error {
	errors := err.(validator.ValidationErrors)
	return errValidate[errors[0].Tag()]
}
