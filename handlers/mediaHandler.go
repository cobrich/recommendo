package handlers

import (
	"net/http"

	"github.com/cobrich/recommendo/service"
)

type MediaHandler struct {
	s *service.MediaService
}

func NewMediaHandler(s *service.MediaService) *MediaHandler {
	return &MediaHandler{s: s}
}

func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {

}