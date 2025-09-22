package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cobrich/recommendo/models"
)

type UserRepo struct {
	db DBTX
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) WithTx(tx *sql.Tx) *UserRepo {
	return &UserRepo{db: tx}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var createdUser models.User

	query := `
		INSERT INTO users (user_name, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING user_id, user_name, email, created_at`

	err := r.db.QueryRowContext(ctx, query, user.UserName, user.Email, user.PasswordHash).Scan(
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

func (r *UserRepo) GetUsers(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	// 1. Gettig total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM users"
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Если записей нет, нет смысла делать второй запрос
	if total == 0 {
		return []models.User{}, 0, nil
	}

	// 2. Getting page datas
	offset := (page - 1) * limit

	query := "SELECT user_id, user_name, created_at FROM users ORDER BY user_name LIMIT $1 OFFSET $2"

	var users []models.User

	sqlRows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer sqlRows.Close()

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		users = append(users, user)

	}
	return users, total, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	query := "SELECT user_id, user_name, created_at FROM users WHERE user_id = $1"

	var user models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.UserName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return models.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserFriends(ctx context.Context, userID, page, limit int) ([]models.User, int64, error) {

	// 1. Gettig total count
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM
		    follows f1
		JOIN
		    follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		WHERE
		    f1.follower_id = $1;
	`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Если записей нет, нет смысла делать второй запрос
	if total == 0 {
		return []models.User{}, 0, nil
	}

	offset := (page - 1) * limit

	dataQuery := `
		SELECT
		    u.user_id,
		    u.user_name,
		    u.created_at
		FROM
		    follows f1
		JOIN
		    follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		JOIN
		    users u ON u.user_id = f1.following_id
		WHERE
		    f1.follower_id = $1;
		ORDER BY
		    u.user_name -- <-- ВАЖНО: Пагинация без сортировки не имеет смысла!
		LIMIT $2 OFFSET $3; -- <-- Новые параметры
	`

	sqlRows, err := r.db.QueryContext(ctx, dataQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user friends: %w", err)
	}

	defer sqlRows.Close()

	var users []models.User

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, total, nil
}

func (r *UserRepo) GetUserFollowers(ctx context.Context, userID, page, limit int) ([]models.User, int64, error) {
	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM
		    follows f
		WHERE
		    f.following_id = $1;
	`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if total == 0 {
		return []models.User{}, 0, nil
	}

	offset := (page - 1) * limit

	dataQuery := `
		SELECT
		    u.user_id,
		    u.user_name,
		    u.created_at
		FROM
		    follows f
		JOIN
		    users u ON u.user_id = f.follower_id
		WHERE
		    f.following_id = $1;
		ORDER BY
			u.user_name 
		LIMIT $2 OFFSET $3;
	`

	sqlRows, err := r.db.QueryContext(ctx, dataQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user friends: %w", err)
	}

	defer sqlRows.Close()

	var users []models.User

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, total, nil
}

func (r *UserRepo) GetUserFollowings(ctx context.Context, userID, page, limit int) ([]models.User, int64, error) {

	// 1
	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM
		    follows f
		WHERE
		    f.follower_id = $1;
	`

	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if total == 0 {
		return []models.User{}, 0, nil
	}

	// 2.
	offset := (page - 1) * limit

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
		ORDER BY
			u.user_name
		LIMIT $2 OFFSET $3;
	`

	sqlRows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user friends: %w", err)
	}

	for sqlRows.Next() {
		var user models.User

		if err := sqlRows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}
	return users, total, nil
}

// FindUserByEmail ищет пользователя по email. Возвращает хеш пароля для проверки в сервисе.
func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := "SELECT user_id, user_name, email, password_hash, created_at FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, err // err может быть sql.ErrNoRows, это нормально
	}
	return user, nil
}

// FindUserByIDWithPassword получает ВСЕ данные пользователя, включая хеш.
func (r *UserRepo) FindUserByIDWithPassword(ctx context.Context, id int) (models.User, error) {
	var user models.User
	// Этот запрос выбирает все поля, включая password_hash
	query := "SELECT user_id, user_name, email, password_hash, created_at FROM users WHERE user_id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, err // sql.ErrNoRows будет обработан в сервисе
	}
	return user, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, userID int) error {
	query := "DELETE FROM users WHERE user_id = $1"

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, userID int, userName string) (models.User, error) {
	var user models.User

	query := "UPDATE users SET user_name = $1 WHERE user_id = $2 RETURNING user_id, user_name, email, created_at"

	err := r.db.QueryRowContext(ctx, query, userName, userID).Scan(&user.ID, &user.UserName, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, sql.ErrNoRows
		}
		return models.User{}, fmt.Errorf("failed to update user name")
	}

	return user, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID int, newPasswordHash []byte) error {
	query := "UPDATE users SET password_hash = $1 WHERE user_id = $2"
	result, err := r.db.ExecContext(ctx, query, newPasswordHash, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
