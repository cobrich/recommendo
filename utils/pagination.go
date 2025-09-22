package utils

import (
	"errors"
	"net/http"
	"strconv"
)

// (вспомогательная функция для парсинга, чтобы не дублировать код)
func ParsePaginationParams(r *http.Request) (page, limit int, err error) {
	// Получаем параметры из URL, например /users?page=2&limit=25
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// --- Устанавливаем значения по умолчанию ---
	page = 1
	limit = 20 // Хорошее значение по умолчанию

	// --- Парсим и валидируем ---
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, errors.New("invalid 'page' parameter: must be a positive integer")
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, errors.New("invalid 'limit' parameter: must be a positive integer")
		}
	}

	// (Опционально) Ограничиваем максимальный limit, чтобы защититься от DoS-атак
	if limit > 100 {
		limit = 100
	}

	return page, limit, nil
}
