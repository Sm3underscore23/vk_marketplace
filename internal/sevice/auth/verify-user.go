package auth

import (
	"context"
	"errors"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (a *authService) VerifyUser(ctx context.Context, userData models.UserData) (int, error) {
	dataForVerify, err := a.userRepo.VerifyUser(ctx, userData.Login)
	if errors.Is(err, models.ErrorWrongLoginOrPassword) {
		return dataForVerify.UserID, models.ErrorWrongLoginOrPassword
	}
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.VerifyUser",
			logger.Error, err,
		)
		return dataForVerify.UserID, models.ErrorDb
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dataForVerify.PasswordHash), []byte(userData.Password)); err != nil {
		return dataForVerify.UserID, models.ErrorWrongLoginOrPassword
	}

	return dataForVerify.UserID, nil
}
