package repo

import (
	"context"
	"database/sql"
	"fmt"
)

type FollowRepo struct {
	db DBTX
}

func NewFollowRepo(db *sql.DB) *FollowRepo {
	return &FollowRepo{db: db}
}

func (r *FollowRepo) WithTx(tx *sql.Tx) *FollowRepo {
    return &FollowRepo{db: tx}
}


func (r *FollowRepo) CreateFollow(ctx context.Context, followerID, followingID int) error {
	query := `
        INSERT INTO follows (follower_id, following_id)
        VALUES ($1, $2)
		`

	result, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		// Если произошла ошибка (например, нарушение UNIQUE constraint), мы ее получим.
		return fmt.Errorf("failed to create following: %w", err)
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

func (r *FollowRepo) DeleteFollow(ctx context.Context, followerID, followingID int) error {
	query := `
		DELETE FROM follows 
		WHERE follower_id = $1 AND following_id= $2
		`

	result, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		// Если произошла ошибка (например, нарушение UNIQUE constraint), мы ее получим.
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	// 3. (Опционально, но хорошая практика) Проверяем, что была затронута ровно одна строка.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Маловероятная ошибка, но проверка не помешает.
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *FollowRepo) AreUsersFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM follows f1
			JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
			WHERE f1.follower_id = $1 AND f1.following_id = $2
		)
	`
	var areFriends bool

	err := r.db.QueryRowContext(ctx, query, userID1, userID2).Scan(&areFriends)
	if err != nil {
		return false, fmt.Errorf("failed to check friendship: %w", err)
	}

	return areFriends, nil
}

func (r *FollowRepo) DeleteAllUserFollows(ctx context.Context, userID int) error {
    query := "DELETE FROM follows WHERE follower_id = $1 OR following_id = $1"
    
    _, err := r.db.ExecContext(ctx, query, userID)
    if err != nil {
        return fmt.Errorf("failed to delete all user follows: %w", err)
    }
    return nil
}