package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/cobrich/recommendo/repo"
)

var ErrFollowNotFound = errors.New("follow relationship not found")

type FollowService struct {
	r      *repo.FollowRepo
	logger *slog.Logger
}

func NewFollowService(r *repo.FollowRepo, logger *slog.Logger) *FollowService {
	return &FollowService{r: r, logger: logger}
}

func (s *FollowService) CreateFollow(ctx context.Context, fromId, toID int) error {
	err := s.r.CreateFollow(ctx, fromId, toID)
	if err != nil {
		return err
	}
	return nil
}

func (s *FollowService) DeleteFollow(ctx context.Context, fromId, toID int) error {
	err := s.r.DeleteFollow(ctx, fromId, toID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// И переводит ее в понятную для хендлера ошибку бизнес-логики
			return ErrFollowNotFound
		}
		// Все остальные ошибки пробрасываем как есть (это могут быть ошибки БД)
		return err
	}
	return nil
}

func (s *FollowService) AreUsersFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	areFriends, err := s.r.AreUsersFriends(ctx, userID1, userID2)
	if err != nil {
		return false, err
	}
	return areFriends, nil
}
