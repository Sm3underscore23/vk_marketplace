package models

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"errors"`
}

var (
	ErrParseConfig = fmt.Errorf("failed to parse config")

	ErrorJWTParse          = fmt.Errorf("failed to parse JWT")
	ErrorInvalidAuthHeader = fmt.Errorf("jws: invalid or empty Authorization header")
	ErrorSingingMethod     = fmt.Errorf("jws: unexpected signing method")
	ErrorTokenClaims       = fmt.Errorf("jws: invalid or expired token")
	ErrorUserIDClaims      = fmt.Errorf("jws: user_id not found or invalid in tokens")
	ErrorUserLoginClaims   = fmt.Errorf("jws:user_login not found or invalid in token")

	ErrorGetUserIDCtx    = fmt.Errorf("ctx: user_id not found")
	ErrorGetUserLoginctx = fmt.Errorf("ctx:user_login not found")

	ErrorInvalidReqBody = fmt.Errorf("invalid request body")

	ErrorLoginAlreadyExists   = fmt.Errorf("such login already exists")
	ErrorWrongLoginOrPassword = fmt.Errorf("wrong login or password")

	ErrorLoginValidation                 = fmt.Errorf("login must contain only letters, numbers, and underscores, and be at least 4 characters long, max 16 characters long")
	ErrorPasswordValidation              = fmt.Errorf("password must contain only letters, numbers, and be at least 4 letters, 1 digit and 1 symbol, max 64 characters long")
	ErrorTitleValidation                 = fmt.Errorf("title must be contain only letters, numbers, and spaces and be at least 3 characters long, max 32 characters long")
	ErrorDescriptionValidation           = fmt.Errorf("description too long, max 1000 characters long")
	ErrorPriceValidation                 = fmt.Errorf("too high/low price or contains more then 2 decimal places, max 9 999 999.99")
	ErrorImageURLValidation              = fmt.Errorf("image URL must be a valid URL, accessible, in a supported image format (e.g., JPG, PNG, GIF), and not exceed 5Mb in size")
	ErrorImageInvalidURL                 = fmt.Errorf("invalid image URL")
	ErrorImageNotAccessible              = fmt.Errorf("image not accessible")
	ErrorImageNonSuccessStatusCode       = fmt.Errorf("image returned non-2xx status code")
	ErrorImageInvalidContentType         = fmt.Errorf("invalid content type (not an image)")
	ErrorImageInvalidContentLengthHeader = fmt.Errorf("invalid content length header")
	ErrorImageTooLarge                   = fmt.Errorf("image exceeds size limit, max 5Mb")
	ErrorImageGETRequestFailed           = fmt.Errorf("failed to perform GET request for image")
	ErrorImageReadFailed                 = fmt.Errorf("failed to read image data")

	ErrorUnknownURIParam         = fmt.Errorf("unknown URI parametr")
	ErrorInvalidSordByURIParam   = fmt.Errorf("invalid sort_by value")
	ErrorInvalidOrderURIParam    = fmt.Errorf("invalid order value")
	ErrorInvalidPriceMinURIParam = fmt.Errorf("invalid price_min")
	ErrorInvalidPriceMaxURIParam = fmt.Errorf("invalid price_max")
	ErrorInvalidPricesURIParam   = fmt.Errorf("price filter failed: price_min > price_max")

	ErrorInvalidCreatedAtURIParam = fmt.Errorf("invalid created_after, expected DateOnly format")
	ErrorInvalidLimitURIParam     = fmt.Errorf("invalid limit")

	ErrorShortCursor = fmt.Errorf("cursor too short")
	ErrorCursorParse = fmt.Errorf("failed to parse cursor")

	ErrorPasswordHashGenerate = fmt.Errorf("failed to generate password hash")
	ErrorJWTGenerate          = fmt.Errorf("failed to generate JWT token")
	ErrorCursorGenerate       = fmt.Errorf("failed to generate cursor token")

	ErrorDb = fmt.Errorf("db error")

	errorToStatus = map[error]int{
		ErrorJWTParse:        http.StatusUnauthorized,
		ErrorUserLoginClaims: http.StatusUnauthorized,

		ErrorInvalidReqBody: http.StatusBadRequest,

		ErrorLoginAlreadyExists:   http.StatusBadRequest,
		ErrorWrongLoginOrPassword: http.StatusBadRequest,

		ErrorLoginValidation:                 http.StatusBadRequest,
		ErrorPasswordValidation:              http.StatusBadRequest,
		ErrorTitleValidation:                 http.StatusBadRequest,
		ErrorDescriptionValidation:           http.StatusBadRequest,
		ErrorPriceValidation:                 http.StatusBadRequest,
		ErrorImageURLValidation:              http.StatusBadRequest,
		ErrorImageInvalidURL:                 http.StatusBadRequest,
		ErrorImageNotAccessible:              http.StatusBadRequest,
		ErrorImageNonSuccessStatusCode:       http.StatusBadRequest,
		ErrorImageInvalidContentType:         http.StatusBadRequest,
		ErrorImageInvalidContentLengthHeader: http.StatusBadRequest,
		ErrorImageTooLarge:                   http.StatusBadRequest,
		ErrorImageGETRequestFailed:           http.StatusBadRequest,
		ErrorImageReadFailed:                 http.StatusBadRequest,

		ErrorUnknownURIParam:          http.StatusBadRequest,
		ErrorInvalidSordByURIParam:    http.StatusBadRequest,
		ErrorInvalidOrderURIParam:     http.StatusBadRequest,
		ErrorInvalidPriceMinURIParam:  http.StatusBadRequest,
		ErrorInvalidPriceMaxURIParam:  http.StatusBadRequest,
		ErrorInvalidCreatedAtURIParam: http.StatusBadRequest,
		ErrorInvalidLimitURIParam:     http.StatusBadRequest,

		ErrorShortCursor: http.StatusBadRequest,
		ErrorCursorParse: http.StatusBadRequest,
	}
)

func GetStatusCode(err error) int {
	if statusCode, ok := errorToStatus[err]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}
