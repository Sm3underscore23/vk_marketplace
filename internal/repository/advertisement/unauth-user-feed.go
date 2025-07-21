package advertisement

// import (
// 	"context"
// 	"fmt"
// 	"marketplace/internal/models"

// 	sq "github.com/Masterminds/squirrel"
// )

// func (a *adRepository) UnAuthUserFeed(
// 	ctx context.Context,
// 	adID int,
// 	uriParams models.FeedURIParams,
// ) (models.UnAuthUserFeed, error) {
// 	var (
// 		adsForFeed models.UnAuthUserFeed
// 		limit      = uriParams.Limit
// 	)

// 	queryOrderBy := orderByBuilder(uriParams)

// 	builder := sq.Select(
// 		advertisementsTable+"."+adIDColumn,
// 		adTitleColumn,
// 		adDescriptionColumn,
// 		adPriceColumn,
// 		asImageUrlColumn,
// 		userLoginColumn,
// 		adCreatedAtColumn,
// 	).
// 		From(advertisementsTable).
// 		Join(fmt.Sprintf("%s ON %s.%s = %s.%s", usersTable, advertisementsTable, adAuthorIDColumn, usersTable, userIDColumn)).
// 		OrderBy(queryOrderBy).
// 		Limit(limit+1).PlaceholderFormat(sq.Dollar)

// 	if !uriParams.CreatedAfter.IsZero() {
// 		builder = builder.Where(sq.Gt{advertisementsTable + "." + adCreatedAtColumn: uriParams.CreatedAfter})
// 	}

// 	if uriParams.PriceMin > 0 {
// 		builder = builder.Where(sq.GtOrEq{advertisementsTable + "." + adPriceColumn: uriParams.PriceMin})
// 	}

// 	if uriParams.PriceMax > 0 {
// 		builder = builder.Where(sq.LtOrEq{advertisementsTable + "." + adPriceColumn: uriParams.PriceMax})
// 	}

// 	query, args, err := builder.ToSql()
// 	if err != nil {
// 		return adsForFeed, err
// 	}

// 	rows, err := a.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return adsForFeed, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var adData models.UnAuthUserAdForFeed
// 		var id int

// 		err := rows.Scan(
// 			&id,
// 			&adData.Title,
// 			&adData.Description,
// 			&adData.Price,
// 			&adData.ImageUrl,
// 			&adData.AuthorLogin,
// 			&adData.CreatedAt,
// 		)
// 		if err != nil {
// 			return adsForFeed, err
// 		}
// 		if len(adsForFeed.Ads) == int(limit) {
// 			adsForFeed.LastAdID = id
// 			continue
// 		}
// 		adsForFeed.Ads = append(adsForFeed.Ads, adData)
// 	}

// 	return adsForFeed, nil
// }
