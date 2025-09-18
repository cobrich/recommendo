package service

import (
	"context"

	"github.com/cobrich/recommendo/repo"
)

type FollowService struct {
	r *repo.FollowRepo
}

func NewFriendshipService(r *repo.FollowRepo) *FollowService {
	return &FollowService{r: r}
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
		return err
	}
	return nil
}
