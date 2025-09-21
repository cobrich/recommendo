package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug" // Пакет для получения стека вызовов
)

// NewRecoverer создает middleware, которое перехватывает паники.
func NewRecoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// defer гарантирует, что эта функция будет выполнена в самом конце,
			// даже если внутри next.ServeHTTP произойдет паника.
			defer func() {
				// recover() возвращает nil, если паники не было.
				if err := recover(); err != nil {
					// Если паника произошла, err не будет nil.

					// 1. Логируем ошибку с максимальной детализацией.
					// Уровень ERROR, так как это критическая проблема.
					logger.Error(
						"Panic recovered",
						"error", err,
						// Стек вызовов - это самое важное для отладки паники!
						"stack", string(debug.Stack()),
					)

					// 2. Отправляем клиенту безопасный ответ 500.
					// Никогда не отправляйте детали паники клиенту!
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			// Вызываем следующий обработчик в цепочке.
			// Если паника случится здесь, defer выше ее перехватит.
			next.ServeHTTP(w, r)
		})
	}
}