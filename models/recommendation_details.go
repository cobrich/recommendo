package models

import "time"

type RecommendationDetails struct {
	RecommendationID int       `db:"recommendation_id"`
	Media            MediaItem // Вложенная структура для информации о медиа
	User             User      // Вложенная структура для информации о втором пользователе
	CreatedAt        time.Time `db:"created_at"`
}