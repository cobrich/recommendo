package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cobrich/recommendo/service"
	"github.com/go-chi/chi/v5"
)

type CreateRecommendationRequestDTO struct {
	FromUserID int `json:"from_user_id"`
	ToUserID   int `json:"to_user_id"`
	MediaID    int `json:"media_id"`
}

type RecommendationHandler struct {
	s *service.RecommendationService
}

func NewRecommendationHandler(s *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{s: s}
}

func (h *RecommendationHandler) CreateRecommendation(w http.ResponseWriter, r *http.Request) {
	// 1. Создаем переменную для хранения данных из тела запроса
	var reqDTO CreateRecommendationRequestDTO

	// 2. Декодируем JSON из тела запроса в нашу DTO
	err := json.NewDecoder(r.Body).Decode(&reqDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. (Опционально, но рекомендуется) Проводим базовую валидацию
	if reqDTO.FromUserID <= 0 || reqDTO.ToUserID <= 0 || reqDTO.MediaID <= 0 {
		http.Error(w, "User and media IDs must be positive integers", http.StatusBadRequest)
		return
	}

	// 4. Вызываем сервисный метод, передавая ему данные из DTO
	err = h.s.CreateRecommendation(r.Context(), reqDTO.FromUserID, reqDTO.ToUserID, reqDTO.MediaID)
	if err != nil {
		// 5. Умная обработка ошибок от сервиса
		if strings.Contains(err.Error(), "not found") {
			// Если сервис вернул "user not found" или "media not found"
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
		} else if strings.Contains(err.Error(), "not friends") || strings.Contains(err.Error(), "already been recommended") {
			// Если сервис вернул "users are not friends" или "already recommended"
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
		} else {
			// Все остальные ошибки - это проблемы на нашей стороне
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 6. Если все прошло успешно, отправляем ответ 201 Created
	w.WriteHeader(http.StatusCreated)
	// Можно отправить просто статус, или небольшое подтверждение в JSON
	json.NewEncoder(w).Encode(map[string]string{"status": "recommendation created successfully"})
}

func (h *RecommendationHandler) GetUserRecommendations(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
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
