package auth

import (
	"context"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (a *authService) SignUp(ctx context.Context, userData models.UserData) (int, error) {
	var userID int
	isExists, err := a.userRepo.IsLoginExists(ctx, userData.Login)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp",
			logger.Error, err,
		)
		return userID, models.ErrorDb
	}

	if isExists {
		return userID, models.ErrorLoginAlreadyExists
	}

	passwordHash, err := generatePasswordHash(userData.Password)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp",
			logger.Error, err,
		)
		return userID, models.ErrorPasswordHashGenerate
	}

	userID, err = a.userRepo.CreateUser(ctx, userData.Login, passwordHash)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp",
			logger.Error, err,
		)
		return userID, models.ErrorDb
	}

	return userID, nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
