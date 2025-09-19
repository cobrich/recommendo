package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cobrich/recommendo/models"
	"github.com/cobrich/recommendo/repo"
)

type RecommendationService struct {
	// Собственные зависимости (репозитории)
	r         *repo.RecommendationRepo
	mediaRepo *repo.MediaRepo // Допустим, он может сам создавать медиа

	// Зависимости от ДРУГИХ СЕРВИСОВ
	userService   *UserService
	followService *FollowService
}

// Конструктор теперь принимает все нужные зависимости
func NewRecommendationService(rRepo *repo.RecommendationRepo, mRepo *repo.MediaRepo, uService *UserService, fService *FollowService) *RecommendationService {
	return &RecommendationService{
		r:             rRepo,
		mediaRepo:     mRepo,
		userService:   uService,
		followService: fService,
	}
}

func (s *RecommendationService) CreateRecommendation(ctx context.Context, fromID, toID, mediaID int) error {
	// 1. Check existance of users
	_, err := s.userService.GetUserByID(ctx, fromID)
	if err != nil {
		return err
	}
	_, err = s.userService.GetUserByID(ctx, toID)
	if err != nil {
		return err
	}

	// 2. Check existance of media
	_, err = s.mediaRepo.GetMedia(ctx, mediaID)
	if err != nil {
		return err
	}

	// 3. Check if users are friends
	areFriends, err := s.followService.AreUsersFriends(ctx, fromID, toID)
	if err != nil {
		return err
	}
	if !areFriends {
		return fmt.Errorf("users are not friends")
	}

	// 4. Check is it recommandation first time
	err = s.r.GetRecommendation(ctx, fromID, toID, mediaID)
	if err == nil {
		return fmt.Errorf("this media has already been recommended to this user")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("unexpected error when checking for existing recommendation: %w", err)
	}

	// 5. If not exists, and there is no problems create recomm
	return s.r.CreateRecommendation(ctx, fromID, toID, mediaID)
}

func (s *RecommendationService) GetRecommendations (ctx context.Context, userID int, direction string) ([]models.RecommendationDetails, error) {
    if direction == "sent" {
        return s.r.GetSentRecommendations(ctx, userID)
    }
    
    return s.r.GetReceivedRecommendations(ctx, userID)
}