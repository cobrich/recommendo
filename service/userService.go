package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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
	ErrFailedHashPassword = errors.New("failed to hash password")
)

type UserService struct {
	db *sql.DB // <-- Добавляем это поле
	r  *repo.UserRepo
	// Добавляем зависимости от других репозиториев
	followRepo *repo.FollowRepo
	recomRepo  *repo.RecommendationRepo
	logger     *slog.Logger
}

func NewUserService(db *sql.DB, userRepo *repo.UserRepo, followRepo *repo.FollowRepo, recomRepo *repo.RecommendationRepo, logger *slog.Logger) *UserService {
	return &UserService{
		db:         db,
		r:          userRepo,
		followRepo: followRepo,
		recomRepo:  recomRepo,
		logger:     logger,
	}
}

func (s *UserService) Register(ctx context.Context, registerDTO dtos.RegisterUserDTO) (models.User, error) {
	s.logger.Info("Register: Starting user registration", "user_name", registerDTO.UserName, "email", registerDTO.Email)

	// 1. Validation fields
	email, err := utils.CleanAndValidateEmail(registerDTO.Email)
	if err != nil {
		s.logger.Error("Register: Email validation failed", "error", err, "email", registerDTO.Email)
		return models.User{}, err
	}

	isValid, errs := utils.ValidatePassword(registerDTO.Password)
	if !isValid {
		s.logger.Error("Register: Password validation failed", "error", errs)
		return models.User{}, errs
	}

	// 2. Проверка, что пользователь не существует (КРИТИЧЕСКИЙ ШАГ)
	s.logger.Info("Register: Checking if user exists", "email", email)
	_, err = s.r.FindUserByEmail(ctx, email)
	if err == nil {
		s.logger.Warn("Register: User already exists", "email", email)
		return models.User{}, ErrUserExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("Register: Database error while checking user existence", "error", err, "email", email)
		return models.User{}, err
	}
	s.logger.Info("Register: User does not exist, proceeding with creation")

	// 3. Hashing password
	s.logger.Info("Register: Hashing password")
	hashedPassword, err := utils.GetPasswordHash(registerDTO.Password)
	if err != nil {
		s.logger.Error("Register: Password hashing failed", "error", err)
		return models.User{}, err
	}

	// 4. Creating new user object
	userToCreate := models.User{
		UserName:     registerDTO.UserName,
		Email:        email,
		PasswordHash: []byte(hashedPassword), // Убедитесь, что тип совпадает с моделью
	}

	// 5. Сохранение в репозитории
	s.logger.Info("Register: Creating user in database", "user_name", userToCreate.UserName, "email", userToCreate.Email)
	createdUser, err := s.r.CreateUser(ctx, userToCreate)
	if err != nil {
		s.logger.Error("Register: Failed to create user in database", "error", err, "user_name", userToCreate.UserName, "email", userToCreate.Email)
		return models.User{}, err
	}

	s.logger.Info("Register: User created successfully", "user_id", createdUser.ID, "user_name", createdUser.UserName, "email", createdUser.Email)
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

func (s *UserService) GetUsers(ctx context.Context, page, limit int) (*dtos.PaginatedResponseDTO[models.User], error) {
	users, total, err := s.r.GetUsers(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// Вычисляем общее количество страниц
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &dtos.PaginatedResponseDTO[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	user, err := s.r.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserFriends(ctx context.Context, userID, page, limit int) (*dtos.PaginatedResponseDTO[models.User], error) {
	// Вызываем обновленный метод репозитория
	users, total, err := s.r.GetUserFriends(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	// Вычисляем общее количество страниц
	totalPages := 0
	if total > 0 {
		// Формула для вычисления количества страниц
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	// Упаковываем все в DTO для ответа
	return &dtos.PaginatedResponseDTO[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserFollowers(ctx context.Context, userID, page, limit int) (*dtos.PaginatedResponseDTO[models.User], error) {
	users, total, err := s.r.GetUserFollowers(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	// Вычисляем общее количество страниц
	totalPages := 0
	if total > 0 {
		// Формула для вычисления количества страниц
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	// Упаковываем все в DTO для ответа
	return &dtos.PaginatedResponseDTO[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserFollowings(ctx context.Context, userID, page, limit int) (*dtos.PaginatedResponseDTO[models.User], error) {
	users, total, err := s.r.GetUserFollowings(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	// Вычисляем общее количество страниц
	totalPages := 0
	if total > 0 {
		// Формула для вычисления количества страниц
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	// Упаковываем все в DTO для ответа
	return &dtos.PaginatedResponseDTO[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	// 1. Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Failed to begin transaction", "error", err)
		return err
	}
	// defer с recover гарантирует, что если что-то пойдет не так, транзакция будет отменена
	defer tx.Rollback()

	// 2. Создаем экземпляры репозиториев, работающие ВНУТРИ этой транзакции
	userRepoTx := s.r.WithTx(tx)
	followRepoTx := s.followRepo.WithTx(tx)
	recomRepoTx := s.recomRepo.WithTx(tx)

	// 3. Выполняем операции в правильном порядке (от зависимых к основной)
	if err := recomRepoTx.DeleteAllUserRecommendations(ctx, userID); err != nil {
		s.logger.Error("Failed to delete user recommendations", "error", err, "userID", userID)
		return err
	}

	if err := followRepoTx.DeleteAllUserFollows(ctx, userID); err != nil {
		s.logger.Error("Failed to delete user follows", "error", err, "userID", userID)
		return err
	}

	if err := userRepoTx.DeleteUser(ctx, userID); err != nil {
		s.logger.Error("Failed to delete user", "error", err, "userID", userID)
		return err
	}

	// 4. Если все прошло успешно, коммитим транзакцию
	return tx.Commit()
}

func (s *UserService) UpadeUser(ctx context.Context, userID int, userName string) (models.User, error) {
	updatedUser, err := s.r.UpdateUser(ctx, userID, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}
	return updatedUser, nil
}

func (s *UserService) ChangeCurrentUserPassword(ctx context.Context, userID int, changePasswordDto dtos.ChangePasswordDTO) error {
	// 1. Get user by id
	user, err := s.r.FindUserByIDWithPassword(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 2. Compare with current password
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(changePasswordDto.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// 3, Validate and Hash new password
	isValid, errs := utils.ValidatePassword(changePasswordDto.NewPassword)
	if !isValid {
		return errs
	}

	hashedPassword, err := utils.GetPasswordHash(changePasswordDto.NewPassword)
	if err != nil {
		return err
	}

	if err = s.r.UpdatePassword(ctx, userID, []byte(hashedPassword)); err != nil {
		return err
	}

	return nil
}
