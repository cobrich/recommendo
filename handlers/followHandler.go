package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cobrich/recommendo/dtos"
	"github.com/cobrich/recommendo/middleware"
	"github.com/cobrich/recommendo/service"
)

type UserIDRequestDTO struct {
	UserID int `json:"user_id"`
}

type FollowHandler struct {
	s *service.FollowService
}

func NewFriendshiphandler(s *service.FollowService) *FollowHandler {
	return &FollowHandler{s: s}
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
	json.NewEncoder(w).Encode(dtos.StatusResponseDTO{Status: "following was successful"})
}

func (h *FollowHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {
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

	err = h.s.DeleteFollow(r.Context(), currentUserID, requestBody.ToUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 status
	w.Write([]byte(`{"status": "follow deleted Successfully"}`))
}
