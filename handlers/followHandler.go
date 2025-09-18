package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cobrich/recommendo/service"
)

type CreateFollowRequestDTO struct {
	FromUserID int `json:"follower_id"`
	ToUserID   int `json:"following_id"`
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

func (h *FollowHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateFollowRequestDTO

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

	err = h.s.CreateFollow(r.Context(), requestBody.FromUserID, requestBody.ToUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 status
	w.Write([]byte(`{"status": "following was Successfully"}`))
}

func (h *FollowHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateFollowRequestDTO

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

	err = h.s.DeleteFollow(r.Context(), requestBody.FromUserID, requestBody.ToUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 status
	w.Write([]byte(`{"status": "follow deleted Successfully"}`))
}
