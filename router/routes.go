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
		router.Post("/follows", followHandler.CreateFollow)
		// DELETE /follows - delete following
		router.Delete("/follows", followHandler.DeleteFollow)

		router.Post("/recommendations", recommendationHandler.CreateRecommendation)

		// --- User Routes ---
		router.Get("/users", userHandler.GetUsers)
		router.Get("/users/{userID}", userHandler.GetUserByID)

		// --- Follow/Friendship Routes ---
		router.Get("/users/{userID}/friends", userHandler.GetUserFriends)
		router.Get("/users/{userID}/followers", userHandler.GetUserFollowers)
		router.Get("/users/{userID}/followings", userHandler.GetUserFollowings)

		// --- Recommendation Routes ---
		router.Get("/users/{userID}/recommandations", recommendationHandler.GetUserRecommendations)

		router.Get("/media", mediaHandler.GetMedia)

	})

	return router
}
