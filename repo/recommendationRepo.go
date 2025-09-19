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

// GetSentRecommendations возвращает список рекомендаций, ОТПРАВЛЕННЫХ пользователем.
func (r *RecommendationRepo) GetSentRecommendations(ctx context.Context, userID int) ([]models.RecommendationDetails, error) {
	// SQL-запрос, который объединяет 3 таблицы: recommendations, media_items и users (для получателя).
	query := `
		SELECT
			r.recommendation_id,
			r.created_at,
			
			-- Поля для media_items
			m.media_id, m.item_type, m.name, m.year, m.author, m.created_at,
			
			-- Поля для users (получателя рекомендации)
			u.user_id, u.user_name, u.created_at
		FROM
			recommendations r
		JOIN
			media_items m ON r.media_id = m.media_id
		JOIN
			-- Присоединяем информацию о том, КОМУ порекомендовали
			users u ON r.to_user_id = u.user_id
		WHERE
			r.from_user_id = $1
		ORDER BY
			r.created_at DESC;
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []models.RecommendationDetails
	for rows.Next() {
		var rec models.RecommendationDetails
		// Сканируем результат в поля нашей "богатой" структуры.
		// Обратите внимание на вложенные поля rec.Media и rec.User.
		if err := rows.Scan(
			&rec.RecommendationID,
			&rec.CreatedAt,
			&rec.Media.ID, &rec.Media.Type, &rec.Media.Name, &rec.Media.Year, &rec.Media.Author, &rec.Media.CreatedAt,
			&rec.User.ID, &rec.User.UserName, &rec.User.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sent recommendation row: %w", err)
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}


// GetReceivedRecommendations возвращает список рекомендаций, ПОЛУЧЕННЫХ пользователем.
func (r *RecommendationRepo) GetReceivedRecommendations(ctx context.Context, userID int) ([]models.RecommendationDetails, error) {
	// Запрос очень похож, но меняются условия в JOIN и WHERE.
	query := `
		SELECT
			r.recommendation_id,
			r.created_at,
			
			m.media_id, m.item_type, m.name, m.year, m.author, m.created_at,
			
			-- Присоединяем информацию о том, КТО порекомендовал
			u.user_id, u.user_name, u.created_at
		FROM
			recommendations r
		JOIN
			media_items m ON r.media_id = m.media_id
		JOIN
			-- Присоединяем информацию об ОТПРАВИТЕЛЕ рекомендации
			users u ON r.from_user_id = u.user_id
		WHERE
			r.to_user_id = $1 -- <-- Главное отличие здесь
		ORDER BY
			r.created_at DESC;
	`
    
    // Код для выполнения запроса и сканирования будет точно таким же, как в GetSentRecommendations
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get received recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []models.RecommendationDetails
	for rows.Next() {
		var rec models.RecommendationDetails
		if err := rows.Scan(
			&rec.RecommendationID,
			&rec.CreatedAt,
			&rec.Media.ID, &rec.Media.Type, &rec.Media.Name, &rec.Media.Year, &rec.Media.Author, &rec.Media.CreatedAt,
			&rec.User.ID, &rec.User.UserName, &rec.User.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan received recommendation row: %w", err)
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}