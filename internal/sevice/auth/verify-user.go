package auth

import (
	"context"
	"errors"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (a *authService) VerifyUser(ctx context.Context, userData models.UserData) (string, error) {
	var tokenJWT string
	dataForVerify, err := a.userRepo.VerifyUser(ctx, userData.Login)
	if errors.Is(err, models.ErrorWrongLoginOrPassword) {
		return tokenJWT, models.ErrorWrongLoginOrPassword
	}
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.VerifyUser",
			logger.Error, err,
		)
		return tokenJWT, models.ErrorDb
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dataForVerify.PasswordHash), []byte(userData.Password)); err != nil {
		return tokenJWT, models.ErrorWrongLoginOrPassword
	}

	tokenJWT, err = a.generateJWT(ctx, dataForVerify.UserID, userData.Login)
	if err != nil {
		return tokenJWT, err
	}

	return tokenJWT, nil
}
