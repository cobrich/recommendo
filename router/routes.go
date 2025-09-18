package router

import (
	"net/http"

	"github.com/cobrich/recommendo/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *handlers.UserHandler, followHandler *handlers.FollowHandler) http.Handler {
	router := chi.NewRouter()

	// --- User Routes ---
	router.Get("/users", userHandler.GetUsers)
	router.Get("/users/{userID}", userHandler.GetUserByID) // <-- изменил на /users

	// --- Follow/Friendship Routes ---
	router.Get("/users/{userID}/friends", userHandler.GetUserFriends)
	router.Get("/users/{userID}/followers", userHandler.GetUserFollowers)
	router.Get("/users/{userID}/followings", userHandler.GetUserFollowings)

	// POST /follows - create following
	router.Post("/follows", followHandler.CreateFollow)
	// DELETE /follows - delete following
	router.Delete("/follows", followHandler.DeleteFollow)

	// router.Post("/recommendations", recommendoHandler.CreateRecommendation)

	return router
}
