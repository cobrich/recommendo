package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cobrich/recommendo/service"
)

type MediaHandler struct {
	s      *service.MediaService
	logger *slog.Logger
}

func NewMediaHandler(s *service.MediaService, logger *slog.Logger) *MediaHandler {
	return &MediaHandler{s: s, logger: logger}
}

func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	mediaType := queryParams.Get("type")
	mediaName := queryParams.Get("name")

	if mediaType == "" && mediaName == "" {
		http.Error(w, "Please provide at least one search parameter (type or name)", http.StatusBadRequest)
		return
	}

	media_items, err := h.s.FindMedia(r.Context(), mediaType, mediaName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(media_items) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(media_items); err != nil {
		http.Error(w, "Failed to encode media items to JSON", http.StatusInternalServerError)
		return
	}
}
