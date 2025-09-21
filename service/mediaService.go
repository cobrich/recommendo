package service

import (
	"context"
	"log/slog"

	"github.com/cobrich/recommendo/models"
	"github.com/cobrich/recommendo/repo"
)

type MediaService struct {
	r      *repo.MediaRepo
	logger *slog.Logger
}

func NewMediaService(r *repo.MediaRepo, logger *slog.Logger) *MediaService {
	return &MediaService{r: r, logger: logger}
}

func (s *MediaService) FindMedia(ctx context.Context, mtype, name string) ([]models.MediaItem, error) {
	media_items, err := s.r.FindMedia(ctx, mtype, name)
	if err != nil {
		return nil, err
	}
	return media_items, nil
}
