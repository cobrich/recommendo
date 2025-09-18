package models

import "time"

type User struct {
	ID        int       `db:"user_id"`
	UserName      string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}