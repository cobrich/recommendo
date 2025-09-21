package dtos

import "time"

type UserResponseDTO struct {
	ID        int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
