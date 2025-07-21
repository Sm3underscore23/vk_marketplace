package advertisement

import (
	"context"
	"marketplace/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (a *adRepository) CreateAd(ctx context.Context, adData models.AdData) error {
	builder := sq.Insert(advertisementsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(
			adTitleColumn,
			adDescriptionColumn,
			asImageUrlColumn,
			adPriceColumn,
			adAuthorIDColumn,
			adCreatedAtColumn,
		).Values(
		adData.Title,
		adData.Description,
		adData.ImageUrl,
		adData.Price,
		adData.AuthorID,
		adData.CreatedAt,
	)

	query, args := builder.MustSql()

	_, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return  nil
}
