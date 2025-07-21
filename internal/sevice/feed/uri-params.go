package feed

import (
	"context"
	"fmt"
	"marketplace/internal/models"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func generateFeedURI(ctx context.Context, baseURl string, uriParams models.FeedURIParams, cursor string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s/feed?%s=%s&%s=%s&%s=%s",
		baseURl,
		models.URIParamCursor, cursor,
		models.URIParamSortBy, uriParams.SortBy,
		models.URIParamOrder, uriParams.Order,
	))

	if uriParams.PriceMin > 0 {
		builder.WriteString(fmt.Sprintf("&%s=%v",
			models.URIParamPriceMin, uriParams.PriceMin))
	}

	if uriParams.PriceMax > 0 {
		builder.WriteString(fmt.Sprintf("&%s=%v",
			models.URIParamPriceMax, uriParams.PriceMax))
	}

	if !uriParams.CreatedAfter.IsZero() {
		builder.WriteString(fmt.Sprintf("&%s=%T",
			models.URIParamCreatedAfter, uriParams.CreatedAfter))
	}

	builder.WriteString(fmt.Sprintf("&%s=%v",
		models.URIParamLimit, uriParams.Limit))

	return builder.String()
}

func (f *feedService) parseURIParams(query url.Values) (models.FeedURIParams, error) {
	uriParams := models.FeedURIParams{
		SortBy: "date",
		Order:  "desc",
		Limit:  f.defaultLimit,
	}

	for key := range query {
		if ok := models.CheckParams(key); !ok {
			return uriParams, models.ErrorUnknownURIParam
		}
	}

	if sortBy := query.Get(models.URIParamSortBy); sortBy != "" {
		if sortBy != "price" && sortBy != "date" {
			return uriParams, models.ErrorInvalidSordByURIParam
		}
		uriParams.SortBy = sortBy
	}

	if order := query.Get(models.URIParamOrder); order != "" {
		if order != "asc" && order != "desc" {
			return uriParams, models.ErrorInvalidOrderURIParam
		}
		uriParams.Order = order
	}

	if priceMin := query.Get(models.URIParamPriceMin); priceMin != "" {
		priceMin, err := strconv.Atoi(priceMin)
		if err != nil || priceMin < 0 {
			return uriParams, models.ErrorInvalidPriceMinURIParam
		}
		uriParams.PriceMin = priceMin
	}

	if priceMax := query.Get(models.URIParamPriceMax); priceMax != "" {
		priceMax, err := strconv.Atoi(priceMax)
		if err != nil || priceMax > 10_000_000_000 {
			return uriParams, models.ErrorInvalidPriceMaxURIParam
		}
		uriParams.PriceMax = priceMax
	}

	if uriParams.PriceMin > uriParams.PriceMax {
		return uriParams, models.ErrorInvalidPricesURIParam
	}

	if createdAfter := query.Get(models.URIParamCreatedAfter); createdAfter != "" {
		createdAfter, err := time.Parse(time.DateOnly, createdAfter)
		if err != nil {
			return uriParams, models.ErrorInvalidCreatedAtURIParam
		}
		uriParams.CreatedAfter = createdAfter
	}

	if rowLimit := query.Get(models.URIParamLimit); rowLimit != "" {
		limit, err := strconv.ParseUint(rowLimit, 10, 64)
		if err != nil || limit <= 0 {
			return uriParams, models.ErrorInvalidLimitURIParam
		}
		uriParams.Limit = limit
	}

	return uriParams, nil
}
