package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cobrich/recommendo/service"
)

type SendFriendRequestDTO struct {
	FromUserID int `json:"from_user_id"`
	ToUserID   int `json:"to_user_id"`
}

type UserIDRequestDTO struct {
	UserID int `json:"user_id"`
}

type FollowHandler struct {
	s *service.FollowService
}

func NewFriendshiphandler(s *service.FollowService) *FollowHandler {
	return &FollowHandler{s: s}
}

func (h *FollowHandler) SendFriendhsip(w http.ResponseWriter, r *http.Request) {
	var requestBody SendFriendRequestDTO

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.FromUserID <= 0 || requestBody.ToUserID <= 0 {
		http.Error(w, "User IDs must be positive integers", http.StatusBadRequest)
		return
	}

	err = h.s.SendFriendRequest(r.Context(), requestBody.FromUserID, requestBody.ToUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 status
	w.Write([]byte(`{"status": "friendship request sent"}`))
}

func (h *FollowHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {}

func (h *FollowHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {}
