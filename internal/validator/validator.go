package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type validatorWrapper struct {
	validator         *validator.Validate
	maxImageSizeBytes int64
}

type CustomValidator interface {
	CheckImage(ctx context.Context, imageURL string) error
	Struct(s interface{}) error
}

func (v *validatorWrapper) Struct(s interface{}) error {
	return v.validator.StructCtx(context.Background(), s)
}

func New(maxImageSizeBytes int64) CustomValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("login", loginValidation)
	v.RegisterValidation("password", passwordValidator)

	v.RegisterValidation("title", titleValidation)
	v.RegisterValidation("description", descriptionValidation)
	v.RegisterValidation("price", priceValidation)
	v.RegisterValidation("image_url", imageValidation)

	return &validatorWrapper{
		validator:         v,
		maxImageSizeBytes: maxImageSizeBytes,
	}
}
