package models

import "time"

type MediaItem struct {
	ID        int       `db:"media_id"`
	Type      MediaType `db:"item_type"`
	Name      string    `db:"name"`
	Year      int       `db:"year"`
	Author    string    `db:"author"`
	CreatedAt time.Time `db:"created_at"`
}