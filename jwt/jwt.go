package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/joho/godotenv"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// --- 1. Функция генерации токена ---
// Вызывается после успешной аутентификации пользователя (проверки логина/пароля).
func GenerateToken(userID int) (string, error) {
	// Устанавливаем время жизни токена, например, 24 часа.
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создаем "заявки" (claims), включая ID пользователя и время истечения.
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Создаем новый токен, указывая алгоритм подписи и "заявки".
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken проверяет подпись токена и возвращает claims в случае успеха.
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Убеждаемся, что алгоритм подписи тот, который мы ожидаем (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err // Ошибка может быть из-за истекшего срока или неверной подписи
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
