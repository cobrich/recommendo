package router

import (
	"log/slog"
	"net/http"

	"github.com/cobrich/recommendo/handlers"
	"github.com/cobrich/recommendo/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(userHandler *handlers.UserHandler, followHandler *handlers.FollowHandler,
	mediaHandler *handlers.MediaHandler, recommendationHandler *handlers.RecommendationHandler, logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// Укажите, с какого источника разрешены запросы.
		// Для разработки идеально подходит адрес вашего Vite dev-сервера.
		AllowedOrigins: []string{"http://localhost:5173"},
		// Разрешенные методы
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		// Разрешенные заголовки
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// Разрешаем отправку cookies (если понадобится в будущем)
		AllowCredentials: true,
		// Время жизни preflight-запроса в секундах
		MaxAge: 300,
	}))

	router.Use(middleware.NewRecoverer(logger))

	router.Use(middleware.NewLogger(logger))

	// Auth Routes
	router.Post("/register", userHandler.RegisterUser)
	router.Post("/login", userHandler.LoginUser)
	router.Get("/users", userHandler.GetUsers)
	router.Get("/users/{userID}", userHandler.GetUserByID)
	router.Get("/users/{userID}/followers", userHandler.GetUserFollowers)
	router.Get("/users/{userID}/followings", userHandler.GetUserFollowings)
	router.Get("/users/{userID}/friends", userHandler.GetUserFriends)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator)

		// POST /follows - create following
		r.Post("/follows", followHandler.CreateFollow)
		// DELETE /follows - delete following
		r.Delete("/follows/{targetUserID}", followHandler.DeleteMyFollow)

		r.Post("/recommendations", recommendationHandler.CreateRecommendation)

		// --- User Routes ---
		r.Get("/me", userHandler.GetCurrentUser)
		r.Delete("/me", userHandler.DeleteCurrentUser)
		r.Patch("/me", userHandler.UpdateCurrentUser)
		r.Put("/me/password", userHandler.ChangeCurrentUserPassword)

		// --- Follow/Friendship Routes ---
		r.Get("/me/friends", userHandler.GetCurrentUserFriends)

		r.Get("/me/followers", userHandler.GetCurrentUserFollowers)
		r.Delete("/me/followers/{targetUserID}", followHandler.DeleteMeFollow)

		r.Get("/me/followings", userHandler.GetCurrentUserFollowings)

		// --- Recommendation Routes ---
		r.Get("/me/recommendations", recommendationHandler.GetCurrentUserRecommendations)
		r.Get("/users/{userID}/recommendations", recommendationHandler.GetUserRecommendations)
		r.Delete("/me/recommendations/{recommendation_id}", recommendationHandler.DeleteRecommendation)

		r.Get("/media", mediaHandler.GetMedia)

	})

	return router
}
