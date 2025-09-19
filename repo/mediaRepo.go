package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cobrich/recommendo/models"
)

type MediaRepo struct {
	DB *sql.DB
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{DB: db}
}

func (r *MediaRepo) FindMedia(ctx context.Context, mtype, name string) ([]models.MediaItem, error) {
	// 1. Начинаем с базового запроса
	query := "SELECT media_id, item_type, name, year, author, created_at FROM media_items WHERE 1=1"

	// 2. Создаем срез для хранения аргументов для плейсхолдеров
	var args []interface{}

	// 3. Динамически добавляем условия в WHERE
	if mtype != "" {
		// Добавляем условие в запрос
		query += fmt.Sprintf(" AND item_type = $%d", len(args)+1)
		// Добавляем значение в срез аргументов
		args = append(args, mtype)
	}

	if name != "" {
		// Используем ILIKE для регистронезависимого поиска по части строки
		query += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		// Добавляем '%' для поиска по префиксу
		args = append(args, "%"+name+"%")
	}

	// 4. (Опционально) Добавляем сортировку и ограничение
	query += " ORDER BY name LIMIT 20"

	// 5. Выполняем финальный, собранный запрос
	sqlRows, err := r.DB.QueryContext(ctx, query, args...) // 'args...' - это специальный синтаксис для передачи среза как отдельных аргументов
	if err != nil {
		return nil, fmt.Errorf("failed to get media items: %w", err)
	}
	defer sqlRows.Close()

	// 6. Сканирование результата (остается без изменений)
	var media_items []models.MediaItem
	for sqlRows.Next() {
		var media_item models.MediaItem
		if err := sqlRows.Scan(
			&media_item.ID,
			&media_item.Type,
			&media_item.Name,
			&media_item.Year,
			&media_item.Author,
			&media_item.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		media_items = append(media_items, media_item)
	}

	return media_items, nil
}

func (r *MediaRepo) GetMedia(ctx context.Context, mediaID int) (models.MediaItem, error) {
	query := "SELECT media_id, item_type, name, year, author, created_at FROM media_items WHERE media_id=$1"

	var media_item models.MediaItem
	if err := r.DB.QueryRowContext(ctx, query, mediaID).Scan(
		&media_item.ID,
		&media_item.Type,
		&media_item.Name,
		&media_item.Year,
		&media_item.Author,
		&media_item.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return models.MediaItem{}, fmt.Errorf("media with id %d not found", mediaID)
		}
		return models.MediaItem{}, fmt.Errorf("error while scanning row: %w", err)

	}
	return media_item, nil
}
