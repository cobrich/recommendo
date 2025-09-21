package models

import "time"

type User struct {
	ID           int       `db:"user_id"`
	UserName     string    `db:"user_name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}