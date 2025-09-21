package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cobrich/recommendo/models"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var createdUser models.User

	query := `
		INSERT INTO users (user_name, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING user_id, user_name, email, created_at`

	err := r.DB.QueryRowContext(ctx, query, user.UserName, user.Email, user.PasswordHash).Scan(
		&createdUser.ID,
		&createdUser.UserName,
		&createdUser.Email,
		&createdUser.CreatedAt,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (r *UserRepo) GetUsers(ctx context.Context) ([]models.User, error) {
	query := "SELECT user_id, user_name, created_at FROM users ORDER BY user_name"

	var users []models.User

	sqlRows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer sqlRows.Close()

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		users = append(users, user)

	}
	return users, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	query := "SELECT user_id, user_name, created_at FROM users WHERE user_id = $1"

	var user models.User

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.UserName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return models.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserFriends(ctx context.Context, id int) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT
		    u.user_id,
		    u.user_name,
		    u.created_at
		FROM
		    -- Начинаем с первой "копии" таблицы follows, чтобы найти, на кого подписан наш пользователь ($1).
		    -- Назовем ее f1 (follows 1).
		    follows f1
		JOIN
		    -- Присоединяем вторую "копию" таблицы follows.
		    -- Назовем ее f2 (follows 2).
		    -- Это нужно, чтобы проверить ОБРАТНУЮ подписку.
		    follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		JOIN
		    -- И только теперь присоединяем таблицу users, чтобы получить информацию о друге.
		    users u ON u.user_id = f1.following_id
		WHERE
		    -- Условие: мы ищем подписки, сделанные нашим пользователем ($1).
		    f1.follower_id = $1;
	`

	sqlRows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user friends: %w", err)
	}

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepo) GetUserFollowers(ctx context.Context, id int) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT
		    u.user_id,
		    u.user_name,
		    u.created_at
		FROM
		    follows f
		JOIN
		    users u ON u.user_id = f.follower_id
		WHERE
		    -- Условие: мы ищем тех, кто подписан на user.
		    f.following_id = $1;
	`

	sqlRows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user friends: %w", err)
	}

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepo) GetUserFollowings(ctx context.Context, id int) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT
		    u.user_id,
		    u.user_name,
		    u.created_at
		FROM
		    follows f
		JOIN
		    users u ON u.user_id = f.following_id
		WHERE
		    -- Условие: мы ищем тех, на кого подписан user ($1).
		    f.follower_id = $1;
	`

	sqlRows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user friends: %w", err)
	}

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// FindUserByEmail ищет пользователя по email. Возвращает хеш пароля для проверки в сервисе.
func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := "SELECT user_id, user_name, email, password_hash, created_at FROM users WHERE email = $1"
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, err // err может быть sql.ErrNoRows, это нормально
	}
	return user, nil
}
