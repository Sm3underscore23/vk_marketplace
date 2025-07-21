package user

import (
	"marketplace/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable = "users"

	userIDColumn           = "id"
	userLoginColumn        = "user_login"
	userPasswordHashColumn = "password_hash"
	userCreatedAtColumn    = "created_at"
	userUpdatedAtColumn    = "updated_at"
)

type userRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}
