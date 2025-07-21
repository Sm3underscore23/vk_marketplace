package advertisement

import (
	"context"
	"fmt"
	"marketplace/internal/models"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

func (a *adRepository) Feed(
	ctx context.Context,
	uriParams models.FeedURIParams,
	lastAdData models.LastAdData,
	userLogin string,
) ([]models.AdForFeed, models.LastAdData, error) {
	var (
		adsForFeed    []models.AdForFeed
		newLastAdData models.LastAdData
	)

	isFirstFeed := lastAdData == (models.LastAdData{})

	builder := sq.Select(
		advertisementsTable+"."+adIDColumn,
		adTitleColumn,
		adDescriptionColumn,
		adPriceColumn,
		asImageUrlColumn,
		userLoginColumn,
		adCreatedAtColumn,
	).
		From(advertisementsTable).
		Join(fmt.Sprintf("%s ON %s.%s = %s.%s", usersTable, advertisementsTable, adAuthorIDColumn, usersTable, userIDColumn)).
		PlaceholderFormat(sq.Dollar).
		Limit(uriParams.Limit + 1)

	if !isFirstFeed {
		whereCursor := buildCursorWhere(lastAdData, uriParams.SortBy, uriParams.Order)
		builder = builder.Where(whereCursor)
	}

	orderBy := buildOrderBy(uriParams.SortBy, uriParams.Order)
	builder = builder.OrderBy(orderBy)

	if !uriParams.CreatedAfter.IsZero() {
		builder = builder.Where(sq.Gt{advertisementsTable + "." + adCreatedAtColumn: uriParams.CreatedAfter})
	}

	if uriParams.PriceMin > 0 {
		builder = builder.Where(sq.GtOrEq{advertisementsTable + "." + adPriceColumn: uriParams.PriceMin})
	}
	if uriParams.PriceMax > 0 {
		builder = builder.Where(sq.LtOrEq{advertisementsTable + "." + adPriceColumn: uriParams.PriceMax})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, models.LastAdData{}, err
	}

	rows, err := a.db.Query(ctx, query, args...)
	if err != nil {
		return nil, models.LastAdData{}, err
	}
	defer rows.Close()

	var id int

	for rows.Next() {
		var adData models.AdForFeed

		err := rows.Scan(
			&id,
			&adData.Title,
			&adData.Description,
			&adData.Price,
			&adData.ImageUrl,
			&adData.AuthorLogin,
			&adData.CreatedAt,
		)
		if err != nil {
			return nil, models.LastAdData{}, err
		}

		if userLogin != "" {
			isMine := adData.AuthorLogin == userLogin
			adData.IsMine = &isMine
		}

		adsForFeed = append(adsForFeed, adData)
	}

	if len(adsForFeed)-1 == 0 {
		return adsForFeed, newLastAdData, nil
	}

	newLastAdData = models.LastAdData{
		UserId: id,
		Price:  adsForFeed[len(adsForFeed)-1].Price,
	}

	return adsForFeed[:len(adsForFeed)-1], newLastAdData, nil
}

func buildOrderBy(sortBy string, order string) string {
	if sortBy == "price" {
		return fmt.Sprintf(
			"%s.%s %s, %s.%s desc",
			advertisementsTable, adPriceColumn, order,
			advertisementsTable, adIDColumn,
		)
	}
	return fmt.Sprintf("%s.%s %s", advertisementsTable, adIDColumn, order)
}

func buildCursorWhere(last models.LastAdData, sortBy string, order string) sq.Sqlizer {
	op := ">"
	if strings.ToLower(order) == "desc" {
		op = "<"
	}

	if sortBy == "date" {
		op := ">="
		if strings.ToLower(order) == "desc" {
			op = "<="
		}
		return sq.Expr(
			fmt.Sprintf("(%s.%s) %s (?)",
				advertisementsTable, adIDColumn,
				op,
			),
			last.UserId,
		)
	}

	return sq.Expr(
		fmt.Sprintf("(%s.%s %s ?) OR (%s.%s = ? AND %s.%s <= ?)",
			advertisementsTable, adPriceColumn, op,
			advertisementsTable, adPriceColumn,
			advertisementsTable, adIDColumn,
		),
		last.Price, last.Price, last.UserId,
	)
}
