package auth

import (
	"context"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *authService) generateJWT(ctx context.Context, userID int, userLogin string) (string, error) {
	var token string
	claims := jwt.MapClaims{
		string(models.UserIDKey):    userID,
		string(models.UserLoginKey): userLogin,
		"exp":                       time.Now().Add(time.Hour * 24).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(a.salt))
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.GenerateJWT",
			logger.Error, err,
		)
		return token, models.ErrorJWTGenerate
	}

	return token, nil
}

func (a *authService) ParseJWT(ctx context.Context, authHeader string) (models.ClaimData, error) {
	var claimData models.ClaimData
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
		return claimData, models.ErrorInvalidJWTFormat
	}

	tokenStr := fields[1]

	secret := []byte(a.salt)
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrorSingingMethod
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return claimData, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claimData, models.ErrorTokenClaims
	}

	userIDFloat, ok := claims[string(models.UserIDKey)].(float64)
	if !ok {
		return claimData, models.ErrorUserIDClaims
	}
	userID := int(userIDFloat)

	userLogin, ok := claims[string(models.UserLoginKey)].(string)
	if !ok {
		return claimData, models.ErrorUserLoginClaims
	}

	claimData = models.ClaimData{
		UserId:    userID,
		UserLogin: userLogin,
	}

	return claimData, nil
}
