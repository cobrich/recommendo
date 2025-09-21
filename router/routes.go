package router

import (
	"net/http"

	"github.com/cobrich/recommendo/handlers"
	"github.com/cobrich/recommendo/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *handlers.UserHandler, followHandler *handlers.FollowHandler,
	mediaHandler *handlers.MediaHandler, recommendationHandler *handlers.RecommendationHandler) http.Handler {
	router := chi.NewRouter()

	// Auth Routes
	router.Post("/register", userHandler.RegisterUser)
	router.Post("/login", userHandler.LoginUser)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator)

		// POST /follows - create following
		r.Post("/follows", followHandler.CreateFollow)
		// DELETE /follows - delete following
		r.Delete("/follows/{targetUserID}", followHandler.DeleteFollow)

		r.Post("/recommendations", recommendationHandler.CreateRecommendation)

		// --- User Routes ---
		r.Get("/users", userHandler.GetUsers)
		r.Get("/me", userHandler.GetCurrentUser)

		// --- Follow/Friendship Routes ---
		r.Get("/me/friends", userHandler.GetUserFriends)
		r.Get("/me/followers", userHandler.GetUserFollowers)
		r.Get("/me/followings", userHandler.GetUserFollowings)

		// --- Recommendation Routes ---
		r.Get("/me/recommandations", recommendationHandler.GetUserRecommendations)

		r.Get("/media", mediaHandler.GetMedia)

	})

	return router
}
