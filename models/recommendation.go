package models

import "time"

type Recommendation struct {
	ID           int       `db:"recommendation_id"`
	FromUserID   int       `db:"from_user_id"` 
	ToUserID     int       `db:"to_user_id"`   
	MediaID      int       `db:"media_id"`     
	CreatedAt    time.Time `db:"created_at"`
}