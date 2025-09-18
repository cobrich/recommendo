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

func (s *FollowService) SendFriendRequest(ctx context.Context, fromId, toID int) error {
	err := s.r.SendFriendRequest(ctx, fromId, toID)
	if err != nil {
		return err
	}
	return nil
}
