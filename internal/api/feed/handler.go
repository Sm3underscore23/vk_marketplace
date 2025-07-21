package feed

import (
	service "marketplace/internal/sevice"
)

type Handler struct {
	authService  service.AuthService
	feedService  service.FeedService
	defaultLimit uint64
}

func New(authService service.AuthService, feedService service.FeedService, defaultLimit uint64) *Handler {
	return &Handler{
		authService:  authService,
		feedService:  feedService,
		defaultLimit: defaultLimit,
	}
}
