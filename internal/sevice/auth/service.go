package auth

import (
	"marketplace/internal/repository"
	service "marketplace/internal/sevice"
)

type authService struct {
	userRepo repository.UserRepository
	salt     string
}

func New(userRepo repository.UserRepository, salt string) service.AuthService {
	return &authService{
		userRepo: userRepo,
		salt:     salt,
	}
}
