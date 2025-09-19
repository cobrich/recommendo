package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cobrich/recommendo/models"
)

type RecommendationRepo struct {
	DB *sql.DB
}

func NewRecommendationRepo(db *sql.DB) *RecommendationRepo {
	return &RecommendationRepo{DB: db}
}

func (r *RecommendationRepo) GetRecommendation(ctx context.Context, fromId, toID, mediaID int) error {
	var recommendation models.Recommendation

	query := `SELECT recommendation_id, from_user_id, to_user_id, media_id, created_at FROM recommendations 
	WHERE from_user_id=$1 and to_user_id=$2 and media_id=$3`

	if err := r.DB.QueryRowContext(ctx, query, fromId, toID, mediaID).Scan(
		&recommendation.ID,
		&recommendation.FromUserID,
		&recommendation.ToUserID,
		&recommendation.MediaID,
		&recommendation.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to get recommendation: %w", err)
	}

	return nil
}

func (r *RecommendationRepo) CreateRecommendation(ctx context.Context, fromId, toID, mediaID int) error {
	query := `
        INSERT INTO recommendations (from_user_id, to_user_id, media_id)
        VALUES ($1, $2, $3)
		`

	result, err := r.DB.ExecContext(ctx, query, fromId, toID, mediaID)
	if err != nil {
		// Если произошла ошибка (например, нарушение UNIQUE constraint), мы ее получим.
		return fmt.Errorf("failed to create recommendation: %w", err)
	}

	// 3. (Опционально, но хорошая практика) Проверяем, что была затронута ровно одна строка.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Маловероятная ошибка, но проверка не помешает.
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected, following not created")
	}

	return nil
}
