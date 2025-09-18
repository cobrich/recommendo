package repo

import (
	"context"
	"database/sql"
	"fmt"
)

type FollowRepo struct {
	DB *sql.DB
}

func NewFollowRepo(db *sql.DB) *FollowRepo {
	return &FollowRepo{DB: db}
}

func (r *FollowRepo) CreateFollow(ctx context.Context, fromID, toID int) error {
	query := `
        INSERT INTO follows (user_id_1, user_id_2)
        VALUES ($1, $2)
		`

	result, err := r.DB.ExecContext(ctx, query, fromID, toID)
	if err != nil {
		// Если произошла ошибка (например, нарушение UNIQUE constraint), мы ее получим.
		return fmt.Errorf("failed to create friendship request: %w", err)
	}

	// 3. (Опционально, но хорошая практика) Проверяем, что была затронута ровно одна строка.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Маловероятная ошибка, но проверка не помешает.
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected, friendship request not created")
	}

	return nil
}