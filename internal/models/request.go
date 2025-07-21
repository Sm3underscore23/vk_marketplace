package models

import "time"

type SingUpInRequest struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"password" validate:"required,password"`
}

type CreateAdRequest struct {
	Title       string `json:"title" validate:"required,title"`
	Price       uint64 `json:"price" validate:"price"`
	Description string `json:"description" validate:"description"`
	ImageUrl    string `json:"image_url" validate:"required"`
}

type FeedURIParams struct {
	SortBy       string
	Order        string
	PriceMin     int
	PriceMax     int
	CreatedAfter time.Time
	Limit        uint64
}

type LastAdData struct {
	UserId    int
	Price     uint64
	CreatedAt time.Time
}

var (
	URIParamCursor       = "cursor"
	URIParamSortBy       = "sort_by"
	URIParamOrder        = "order"
	URIParamPriceMin     = "price_min"
	URIParamPriceMax     = "price_max"
	URIParamCreatedAfter = "created_after"
	URIParamLimit        = "limit"

	uriParams = map[string]struct{}{
		URIParamCursor:       {},
		URIParamSortBy:       {},
		URIParamOrder:        {},
		URIParamPriceMin:     {},
		URIParamPriceMax:     {},
		URIParamCreatedAfter: {},
		URIParamLimit:        {},
	}
)

func CheckParams(param string) bool {
	_, ok := uriParams[param]
	return ok
}
