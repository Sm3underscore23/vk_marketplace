package feed

import (
	service "marketplace/internal/sevice"
)

type Handler struct {
	feedService service.FeedService
}

func New(feedService service.FeedService) *Handler {
	return &Handler{
		feedService: feedService,
	}
}
