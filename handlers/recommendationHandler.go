package handlers

import (
	"net/http"

	"github.com/cobrich/recommendo/service"
)

type RecommendationHandler struct {
	s *service.RecommendationService
}

func NewRecommendationHandler(s *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{s: s}
}

func (h *RecommendationHandler) CreateRecommendation(w http.ResponseWriter, r *http.Request) {

}