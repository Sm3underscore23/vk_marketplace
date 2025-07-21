package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	loginValidRegex  = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	loginLetterRegex = regexp.MustCompile(`[a-zA-Z]`)

	passwordLetterRegex = regexp.MustCompile(`([a-zA-Z].*?)+$`)
	passwordDigitRegex  = regexp.MustCompile(`[0-9]`)
	passwordSybolRegex  = regexp.MustCompile(`[!@#\$%^&*(),.?":{}|<>]`)
)

func loginValidation(fl validator.FieldLevel) bool {
	login := fl.Field().String()
	isValidLetters := loginValidRegex.MatchString(login)
	isValidDigitsAndSybols := len(loginLetterRegex.FindAllString(login, -1)) > 3
	max := len([]rune(login)) <= 16
	return isValidLetters && isValidDigitsAndSybols && max
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	letterCheck := passwordLetterRegex.MatchString(password)
	hasDigit := passwordDigitRegex.MatchString(password)
	hasSymbol := passwordSybolRegex.MatchString(password)
	max := len([]rune(password)) <= 64
	return letterCheck && hasDigit && hasSymbol && max
}
