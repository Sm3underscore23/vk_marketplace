package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (u *userRepository) CreateUser(ctx context.Context, login string, passwordHash string) (int, error) {
	var userID int
	builder := sq.Insert(usersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(
			userLoginColumn,
			userPasswordHashColumn,
		).Values(
		login,
		passwordHash,
	).Suffix("RETURNING " + userIDColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return userID, err
	}

	err = u.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
