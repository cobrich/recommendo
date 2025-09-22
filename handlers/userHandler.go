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
	"github.com/cobrich/recommendo/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	s      *service.UserService
	logger *slog.Logger
}

func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{s: userService, logger: logger}
}


// Registering user
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

// Logging in user
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

// Getting users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сервис теперь возвращает готовую DTO
	paginatedResponse, err := h.s.GetUsers(r.Context(), page, limit)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {

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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}
	user, err := h.s.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}
}


// Getting user friends
func (h *UserHandler) GetCurrentUserFriends(w http.ResponseWriter, r *http.Request) {

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	paginatedResponse, err := h.s.GetUserFriends(r.Context(), currentUserID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Could not process request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}

func (h *UserHandler) GetUserFriends(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paginatedResponse, err := h.s.GetUserFriends(r.Context(), userID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Could not process request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}


// Getting user followers
func (h *UserHandler) GetCurrentUserFollowers(w http.ResponseWriter, r *http.Request) {

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	paginatedResponse, err := h.s.GetUserFollowers(r.Context(), currentUserID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Could not process request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}

func (h *UserHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paginatedResponse, err := h.s.GetUserFollowers(r.Context(), userID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get user followers", "error", err, "userID", userID)
			http.Error(w, "failed to get user followers", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}


// Getting user followings
func (h *UserHandler) GetCurrentUserFollowings(w http.ResponseWriter, r *http.Request) {

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		// Эта ошибка не должна происходить, если middleware работает правильно,
		// но проверка - хорошая практика.
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	paginatedResponse, err := h.s.GetUserFollowings(r.Context(), currentUserID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Could not process request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}

func (h *UserHandler) GetUserFollowings(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	page, limit, err := utils.ParsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paginatedResponse, err := h.s.GetUserFollowings(r.Context(), userID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get user followings", "error", err, "userID", userID)
			http.Error(w, "failed to get user followings", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse)
}


// Deleting user
func (h *UserHandler) DeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.s.DeleteUser(r.Context(), currentUserID)
	if err != nil {
		// Здесь уже есть логгер из сервиса, можно добавить еще один в хендлере
		http.Error(w, "Failed to delete user account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// Updating user
func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get current user id
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Get user name from json body
	var user dtos.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}

	if user.UserName == nil {
		http.Error(w, "empty user name", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.s.UpadeUser(r.Context(), currentUserID, *user.UserName)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to upadate user", http.StatusInternalServerError)
		return
	}

	responseDTO := dtos.UserResponseDTO{
		ID:        updatedUser.ID,
		UserName:  updatedUser.UserName,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}


// Changing user password
func (h *UserHandler) ChangeCurrentUserPassword(w http.ResponseWriter, r *http.Request) {
	// 1. Get current user id
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Get passwords from request body
	var changePasswordDto dtos.ChangePasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&changePasswordDto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Validate for empty
	if changePasswordDto.CurrentPassword == "" || changePasswordDto.NewPassword == "" {
		http.Error(w, "Fields 'current_password' and 'new_password' are required", http.StatusBadRequest)
		return
	}

	// 4. Call service
	if err := h.s.ChangeCurrentUserPassword(r.Context(), currentUserID, changePasswordDto); err != nil {
		var passwordValidationErrors utils.PasswordErrors
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			http.Error(w, "Invalid current password", http.StatusForbidden) // 403
		case errors.As(err, &passwordValidationErrors):
			// Если ошибка - это наша структура ошибок валидации
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest) // 400
			json.NewEncoder(w).Encode(passwordValidationErrors)
		default:
			h.logger.Error("Failed to change password", "error", err, "userID", currentUserID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}



