package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cobrich/recommendo/dtos"
	"github.com/cobrich/recommendo/middleware"
	"github.com/cobrich/recommendo/service"
	"github.com/go-chi/chi/v5"
)

type RecommendationHandler struct {
	s      *service.RecommendationService
	logger *slog.Logger
}

func NewRecommendationHandler(s *service.RecommendationService, logger *slog.Logger) *RecommendationHandler {
	return &RecommendationHandler{s: s, logger: logger}
}

func (h *RecommendationHandler) CreateRecommendation(w http.ResponseWriter, r *http.Request) {
	// 1. Создаем переменную для хранения данных из тела запроса
	var reqDTO dtos.CreateRecommendationRequestDTO

	// 2. Декодируем JSON из тела запроса в нашу DTO
	err := json.NewDecoder(r.Body).Decode(&reqDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	// 3. (Опционально, но рекомендуется) Проводим базовую валидацию
	if currentUserID <= 0 || reqDTO.ToUserID <= 0 || reqDTO.MediaID <= 0 {
		http.Error(w, "User and media IDs must be positive integers", http.StatusBadRequest)
		return
	}

	// 4. Вызываем сервисный метод, передавая ему данные из DTO
	err = h.s.CreateRecommendation(r.Context(), currentUserID, reqDTO.ToUserID, reqDTO.MediaID)
	if err != nil {
		// 5. Умная обработка ошибок от сервиса
		switch {
		case errors.Is(err, service.ErrTargetUserNotFound) || errors.Is(err, service.ErrMediaNotFound):
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
			return
		case errors.Is(err, service.ErrNotFriends) || errors.Is(err, service.ErrAlreadyRecommended):
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		default:
			// Все остальные ошибки - это проблемы на нашей стороне
			// log.Printf("Internal server error: %v", err) // Хорошо бы логировать для себя
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// 6. Если все прошло успешно, отправляем ответ 201 Created
	w.WriteHeader(http.StatusCreated)
	// Можно отправить просто статус, или небольшое подтверждение в JSON
	json.NewEncoder(w).Encode(map[string]string{"status": "recommendation created successfully"})
}

func (h *RecommendationHandler) GetCurrentUserRecommendations(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	direction := r.URL.Query().Get("direction")
	recommendations, err := h.s.GetRecommendations(r.Context(), currentUserID, direction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(recommendations) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err = json.NewEncoder(w).Encode(recommendations); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *RecommendationHandler) GetUserRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	
	direction := r.URL.Query().Get("direction")
	recommendations, err := h.s.GetRecommendations(r.Context(), userID, direction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(recommendations) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err = json.NewEncoder(w).Encode(recommendations); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *RecommendationHandler) DeleteRecommendation(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	result := chi.URLParam(r, "recommendation_id")
	recomID, err := strconv.Atoi(result)
	if err != nil {
		h.logger.Warn("Invalid recommendation ID in URL", "error", err, "value", result)
		http.Error(w, "invalid recommendation id", http.StatusBadRequest)
		return
	}
	err = h.s.DeleteRecommendation(r.Context(), currentUserID, recomID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotAuthor):
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, service.ErrUserNotAuthor):
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			h.logger.Error("Failed to delete recommendation", "error", err, "recommendationID", recomID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
