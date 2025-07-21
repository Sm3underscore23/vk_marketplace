package user

import (
	"context"
	"errors"
	"marketplace/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (u *userRepository) VerifyUser(ctx context.Context, login string) (models.UserVerify, error) {
	var userData models.UserVerify
	builder := sq.Select(userIDColumn, userPasswordHashColumn).PlaceholderFormat(sq.Dollar).
		From(usersTable).
		Where(
			sq.Eq{
				userLoginColumn: login,
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return userData, err
	}

	err = u.db.QueryRow(ctx, query, args...).Scan(&userData.UserID, &userData.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return userData, models.ErrorWrongLoginOrPassword
	}
	
	if err != nil {
		return userData, err
	}

	return userData, nil
}