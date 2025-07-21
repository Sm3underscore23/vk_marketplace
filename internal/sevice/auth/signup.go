package auth

import (
	"context"
	"log/slog"
	"marketplace/internal/models"
	"marketplace/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (a *authService) SignUp(ctx context.Context, userData models.UserData) (string, error) {
	var tokenJWT string
	isExists, err := a.userRepo.IsLoginExists(ctx, userData.Login)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp.IsLoginExists",
			logger.Error, err,
		)
		return tokenJWT, models.ErrorDb
	}

	if isExists {
		return tokenJWT, models.ErrorLoginAlreadyExists
	}

	passwordHash, err := generatePasswordHash(userData.Password)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp.generatePasswordHash",
			logger.Error, err,
		)
		return tokenJWT, models.ErrorPasswordHashGenerate
	}

	userID, err := a.userRepo.CreateUser(ctx, userData.Login, passwordHash)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"AuthService.SignUp.CreateUser",
			logger.Error, err,
		)
		return tokenJWT, models.ErrorDb
	}

	tokenJWT, err = a.generateJWT(ctx, userID, userData.Login)
	if err != nil {
		return tokenJWT, err
	}

	return tokenJWT, nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
