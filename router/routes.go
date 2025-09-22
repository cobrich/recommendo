package router

import (
	"log/slog"
	"net/http"

	"github.com/cobrich/recommendo/handlers"
	"github.com/cobrich/recommendo/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *handlers.UserHandler, followHandler *handlers.FollowHandler,
	mediaHandler *handlers.MediaHandler, recommendationHandler *handlers.RecommendationHandler, logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.NewRecoverer(logger))

	router.Use(middleware.NewLogger(logger))

	// Auth Routes
	router.Post("/register", userHandler.RegisterUser)
	router.Post("/login", userHandler.LoginUser)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator)

		// POST /follows - create following
		r.Post("/follows", followHandler.CreateFollow)
		// DELETE /follows - delete following
		r.Delete("/follows/{targetUserID}", followHandler.DeleteMyFollow)

		r.Post("/recommendations", recommendationHandler.CreateRecommendation)

		r.Get("/users", userHandler.GetUsers)

		// --- User Routes ---
		r.Get("/me", userHandler.GetCurrentUser)

		// --- Follow/Friendship Routes ---
		r.Get("/me/friends", userHandler.GetUserFriends)

		r.Get("/me/followers", userHandler.GetUserFollowers)
		r.Delete("/me/followers/{targetUserID}", followHandler.DeleteMeFollow)

		r.Get("/me/followings", userHandler.GetUserFollowings)

		// --- Recommendation Routes ---
		r.Get("/me/recommandations", recommendationHandler.GetUserRecommendations)
		r.Delete("me/recommandations/{recommendation_id}", recommendationHandler.DeleteRecommendation)

		r.Get("/media", mediaHandler.GetMedia)

	})

	return router
}
