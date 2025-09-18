package repo

import (
	"context"
	"database/sql"

	"github.com/cobrich/recommendo/models"
)

type MediaRepo struct {
	DB *sql.DB
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{DB: db}
}

func (r *MediaRepo) GetMedia(ctx context.Context, mtype, name string) ([]models.MediaItem, error) {
	return nil, nil
}