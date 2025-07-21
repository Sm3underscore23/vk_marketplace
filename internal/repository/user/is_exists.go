package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (u *userRepository) IsLoginExists(ctx context.Context, login string) (bool, error) {
	builder := sq.Select(userIDColumn).PlaceholderFormat(sq.Dollar).
		From(usersTable).
		Where(sq.Eq{userLoginColumn: login})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	var id int
	err = u.db.QueryRow(ctx, query, args...).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
