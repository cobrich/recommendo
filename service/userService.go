package service

import (
	"context"

	"github.com/cobrich/recommendo/models"
	"github.com/cobrich/recommendo/repo"
)

type UserService struct {
	r *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{r: userRepo}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.r.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	user, err := s.r.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserFriends(ctx context.Context, id int) ([]models.User, error) {
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	users, err := s.r.GetUserFriends(ctx, id)
	if err != nil {
		return nil, err
	}
	return users, nil
}