package dtos

type CreateRecommendationRequestDTO struct {
	ToUserID   int `json:"to_user_id"`
	MediaID    int `json:"media_id"`
}