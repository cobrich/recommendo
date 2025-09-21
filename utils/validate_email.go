package utils

import (
	"errors"
	"regexp"
	"strings"
)

// Простое регулярное выражение для проверки email.
// В реальности они могут быть гораздо сложнее, но это покрывает 99% случаев.
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// CleanAndValidateEmail обрабатывает и проверяет email.
func CleanAndValidateEmail(email string) (string, error) {
    // 1. Обрезаем пробелы
    email = strings.TrimSpace(email)

    // 2. Приводим к нижнему регистру
    email = strings.ToLower(email)

    // 3. Проверяем формат
    if !emailRegex.MatchString(email) {
        return "", errors.New("некорректный формат email")
    }

    // Если все в порядке, возвращаем очищенный email
    return email, nil
}

// Использование:
// cleanEmail, err := CleanAndValidateEmail("  Test@Example.COM  ")
// if err != nil {
//     // обработать ошибку
// }
// // Теперь cleanEmail == "test@example.com" и его можно сохранять в БД.