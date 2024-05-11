package server

import (
	"net/http"
	"time"

	"github.com/Viet-ph/RSS-Feed-Aggregator/internal/middleware"
	"github.com/Viet-ph/RSS-Feed-Aggregator/internal/service"
)

func NewServer(
	userService *service.UserService,
	feedService *service.FeedService,
	postService *service.PostService,
	feedFollowService *service.FeedFollowService,
) http.Handler {
	go service.StartScraping(feedService, postService, time.Minute)
	mux := http.NewServeMux()
	addRoutes(mux,
		userService,
		feedService,
		postService,
		feedFollowService,
	)
	var handler http.Handler = mux
	handler = middleware.MiddlewareCors(handler)
	return handler
}
