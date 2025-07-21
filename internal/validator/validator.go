package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("login", loginValidation)
	v.RegisterValidation("password", passwordValidator)

	v.RegisterValidation("title", titleValidation)
	v.RegisterValidation("description", descriptionValidation)
	v.RegisterValidation("price", priceValidation)
	v.RegisterValidation("image_url", imageValidation)

	return v
}
