package validator

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	titleValidRegex  = regexp.MustCompile(`^[а-яА-Яa-zA-Z0-9 \/.,\"'<>\(\)\*\:]+$`)
	titleLetterRegex = regexp.MustCompile(`[а-яА-Яa-zA-Z]`)

	allowedImageExtensions = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
		".bmp":  {},
		".webp": {},
		".svg":  {},
	}
)

func titleValidation(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	isValidLetters := titleValidRegex.MatchString(title)
	isValidDigitsAndSybols := len(titleLetterRegex.FindAllString(title, -1)) > 2
	max := len([]rune(title)) <= 32
	return isValidLetters && isValidDigitsAndSybols && max
}

func descriptionValidation(fl validator.FieldLevel) bool {
	desc := fl.Field().String()
	return len(desc) <= 1000
}

func priceValidation(fl validator.FieldLevel) bool {
	price := fl.Field().Uint()
	max := price <= 10_000_000_000
	return max
}

func imageValidation(fl validator.FieldLevel) bool {
	rawURL := fl.Field().String()

	if len([]rune(rawURL)) > 256 {
		return false
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(parsed.Path))
	if ext == "" {
		return false
	}

	if _, ok := allowedImageExtensions[ext]; ok {
		return false
	}

	return false
}
