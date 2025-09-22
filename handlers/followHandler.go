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

type FollowHandler struct {
	s      *service.FollowService
	logger *slog.Logger
}

func NewFriendshiphandler(s *service.FollowService, logger *slog.Logger) *FollowHandler {
	return &FollowHandler{s: s, logger: logger}
}

func (h *FollowHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {
	var requestBody dtos.CreateFollowRequestDTO

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	if currentUserID <= 0 || requestBody.ToUserID <= 0 {
		http.Error(w, "User IDs must be positive integers", http.StatusBadRequest)
		return
	}

	err = h.s.CreateFollow(r.Context(), currentUserID, requestBody.ToUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dtos.StatusResponseDTO{Status: "following was successful"})
}

func (h *FollowHandler) DeleteMyFollow(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	targetUserIDStr := chi.URLParam(r, "targetUserID")
	targetUserID, err := strconv.Atoi(targetUserIDStr)

	if currentUserID <= 0 || targetUserID <= 0 {
		http.Error(w, "User IDs must be positive integers", http.StatusBadRequest)
		return
	}

	err = h.s.DeleteFollow(r.Context(), currentUserID, targetUserID)
	if err != nil {
		if errors.Is(err, service.ErrFollowNotFound) {
			// Если пользователь пытается удалить подписчика, которого нет
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
		} else {
			// Все остальные ошибки - это 500
			h.logger.Error("Failed to remove follower", "error", err, "removerID", currentUserID, "targetID", targetUserID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowHandler) DeleteMeFollow(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	targetUserIDStr := chi.URLParam(r, "targetUserID")
	targetUserID, err := strconv.Atoi(targetUserIDStr)

	if currentUserID <= 0 || targetUserID <= 0 {
		http.Error(w, "User IDs must be positive integers", http.StatusBadRequest)
		return
	}

	err = h.s.DeleteFollow(r.Context(), targetUserID, currentUserID)
	if err != nil {
		if errors.Is(err, service.ErrFollowNotFound) {
			// Если пользователь пытается удалить подписчика, которого нет
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
		} else {
			// Все остальные ошибки - это 500
			h.logger.Error("Failed to remove follower", "error", err, "removerID", currentUserID, "targetID", targetUserID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
