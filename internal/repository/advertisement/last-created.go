package advertisement

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (a *adRepository) LastCreatedAd(ctx context.Context) (int, error) {
	var id int

	builder := sq.Select(adIDColumn).
		From(advertisementsTable).
		OrderBy(adIDColumn + " DESC").Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return id, err
	}

	err = a.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
