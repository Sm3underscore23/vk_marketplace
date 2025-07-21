package models

type ClaimData struct {
	UserId    int
	UserLogin string
}

type FeedResponse struct {
	Ads         []AdForFeed `json:"ads"`
	NextPageURI string      `json:"next_page_uri"`
}
