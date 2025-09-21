package dtos

type CreateFollowRequestDTO struct {
	ToUserID   int `json:"following_id"`
}