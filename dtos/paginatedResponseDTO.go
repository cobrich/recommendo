package dtos

// PaginatedResponseDTO - это универсальная структура для ответов с пагинацией
type PaginatedResponseDTO[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`       // Общее количество записей
	Page       int   `json:"page"`        // Текущая страница
	Limit      int   `json:"limit"`       // Лимит записей на странице
	TotalPages int   `json:"total_pages"` // Общее количество страниц
}