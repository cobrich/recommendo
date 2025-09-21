package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/cobrich/recommendo/dtos"
	"github.com/cobrich/recommendo/middleware"
	"github.com/cobrich/recommendo/service"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{s: userService}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerDTO dtos.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		// Если JSON невалидный - это ошибка клиента
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.s.Register(r.Context(), registerDTO)
	if err != nil {
		// Проверяем тип ошибки из сервиса
		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
		} else {
			// Логируем полную ошибку для себя, а пользователю даем общее сообщение
			// log.Printf("Internal error on user registration: %v", err)
			http.Error(w, "Could not process request", http.StatusInternalServerError)
		}
		return
	}

	responseDTO := dtos.UserResponseDTO{
		ID:        createdUser.ID,
		UserName:  createdUser.UserName,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created для успешного создания ресурса
	json.NewEncoder(w).Encode(responseDTO)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// 1. Getting User Information from Request Body
	var loginDTO dtos.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		// Если JSON невалидный - это ошибка клиента
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 2. Creating token
	token, err := h.s.Login(r.Context(), loginDTO)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			// Все остальные ошибки - это 500
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 3. Send token
	if err := json.NewEncoder(w).Encode(dtos.TokenResponseDTO{Token: token}); err != nil {
		http.Error(w, "Failed to encode token to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.s.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "userID")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	user, err := h.s.GetUserByID(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserFriends(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "userID")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	users, err := h.s.GetUserFriends(r.Context(), currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "userID")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	users, err := h.s.GetUserFollowers(r.Context(), currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserFollowings(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "userID")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	users, err := h.s.GetUserFollowings(r.Context(), currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
	}
}
