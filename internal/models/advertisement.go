package models

import "time"

type AdInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	ImageUrl    string `json:"image_url"`
}

type AdData struct {
	AdInfo
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AdForFeed struct {
	AdInfo
	AuthorLogin string    `json:"author_login"`
	CreatedAt   time.Time `json:"created_at"`
	IsMine      *bool     `json:"is_mine,omitempty"`
}

type Feed struct {
	Ads           []AdForFeed `json:"ads"`
	NextTokenPage string      `json:"next_token_page"`
}
