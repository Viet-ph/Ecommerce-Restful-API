package server

import (
	"net/http"

	"github.com/Viet-ph/RSS-Feed-Aggregator/internal/handler"
	"github.com/Viet-ph/RSS-Feed-Aggregator/internal/middleware"
	"github.com/Viet-ph/RSS-Feed-Aggregator/internal/service"
)

func addRoutes(
	mux *http.ServeMux,
	userService *service.UserService,
	feedService *service.FeedService,
	postService *service.PostService,
	feedFollowService *service.FeedFollowService,
) {
	//Protected routes
	middlewareAuth := middleware.NewMiddlewareAuth(userService)
	mux.Handle("GET /v1/users", middlewareAuth(handler.HandleGetUserByAPIKey()))
	mux.Handle("GET /v1/posts", middlewareAuth(handler.HandleGetPostsByUser(postService)))
	mux.Handle("GET /v1/feed_follows", middlewareAuth(handler.HandleGetFeedFollows(feedFollowService)))
	mux.Handle("POST /v1/feeds", middlewareAuth(handler.HandleCreateFeed(feedService, feedFollowService)))
	mux.Handle("POST /v1/feed_follows", middlewareAuth(handler.HandleCreateFeedFollow(feedFollowService)))
	mux.Handle("DELETE /v1/feed_follows/{feedFollowID}", middlewareAuth(handler.HandleDeleteFeedFollow(feedFollowService)))

	//Unprotected routes
	mux.HandleFunc("GET /api/healthz", handler.Readiness)
	mux.Handle("POST /v1/users", handler.HandleCreateUser(userService))
	mux.Handle("GET /v1/feeds", handler.HandleGetFeeds(feedService))
}
