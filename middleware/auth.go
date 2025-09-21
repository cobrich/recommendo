// file: middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cobrich/recommendo/jwt" // Middleware использует jwt
)

// Определяем кастомный ключ для контекста. Это предотвращает случайные коллизии.
type contextKey string
const UserIDKey contextKey = "userID"

// JWTAuthenticator - это middleware для проверки JWT токена.
func JWTAuthenticator(next http.Handler) http.Handler {
	// http.HandlerFunc - это адаптер, позволяющий использовать обычные функции как http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Проверяем формат "Bearer <token>"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		
		tokenString := headerParts[1]

		// 3. Парсим и валидируем токен с помощью нашего пакета jwt
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. (САМЫЙ ВАЖНЫЙ ШАГ!) Добавляем ID пользователя в контекст запроса.
		// Теперь все последующие хендлеры в цепочке смогут получить этот ID.
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// 5. Вызываем следующий хендлер в цепочке с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает ID пользователя из контекста.
// Возвращает ID и true, если ID найден, иначе 0 и false.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(UserIDKey).(int)
    return userID, ok
}