package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cobrich/recommendo/dtos"
	"github.com/cobrich/recommendo/jwt"
	"github.com/cobrich/recommendo/models"
	"github.com/cobrich/recommendo/repo"
	"github.com/cobrich/recommendo/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials") // Для логина
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	r *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{r: userRepo}
}

func (s *UserService) Register(ctx context.Context, registerDTO dtos.RegisterUserDTO) (models.User, error) {
	// 1. Validation fields
	email, err := utils.CleanAndValidateEmail(registerDTO.Email)
	if err != nil {
		return models.User{}, err

	}
	isValid, errs := utils.ValidatePassword(registerDTO.Password)
	if !isValid {
		return models.User{}, errs
	}

	// 2. Проверка, что пользователь не существует (КРИТИЧЕСКИЙ ШАГ)
	_, err = s.r.FindUserByEmail(ctx, email)
	if err == nil {
		// Ошибки не было, значит пользователь найден
		return models.User{}, ErrUserExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Произошла другая ошибка, не "не найдено"
		return models.User{}, err
	}
	// Если err == sql.ErrNoRows, то все отлично, можно продолжать.

	// 3. Hashing password
	hashedPassword, err := utils.GetPasswordHash(registerDTO.Password)
	if err != nil {
		return models.User{}, err
	}

	// 4. Creating new user object
	userToCreate := models.User{
		UserName:     registerDTO.UserName,
		Email:        email,
		PasswordHash: []byte(hashedPassword), // Убедитесь, что тип совпадает с моделью
	}

	// 5. Сохранение в репозитории
	createdUser, err := s.r.CreateUser(ctx, userToCreate)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, loginDTO dtos.LoginUserDTO) (string, error) {
	// 1. Validate fields for empty
	if loginDTO.Email == "" || loginDTO.Password == "" {
		return "", fmt.Errorf("email and/or password are/is empty")
	}

	// 2. Finding user in db
	user, err := s.r.FindUserByEmail(ctx, loginDTO.Email)
	if err == sql.ErrNoRows {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginDTO.Password)); err != nil {
		return "", ErrInvalidCredentials
	}
	tokenString, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
	users, err := s.r.GetUserFriends(ctx, id)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserFollowers(ctx context.Context, id int) ([]models.User, error) {
	users, err := s.r.GetUserFollowers(ctx, id)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserFollowings(ctx context.Context, id int) ([]models.User, error) {
	users, err := s.r.GetUserFollowings(ctx, id)
	if err != nil {
		return nil, err
	}
	return users, nil
}
