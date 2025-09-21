package utils

import (
	"strings"
	"unicode"
)

// PasswordErrors содержит детальную информацию об ошибках валидации пароля.
type PasswordErrors struct {
	Length     bool `json:"length"`
	HasUpper   bool `json:"has_upper"`
	HasLower   bool `json:"has_lower"`
	HasNumber  bool `json:"has_number"`
	HasSpecial bool `json:"has_special"`
}

// Error возвращает строковое представление ошибок.
// Это делает PasswordErrors совместимым с типом error.
func (e PasswordErrors) Error() string {
	// Собираем сообщения для пользователя
	var messages []string
	if e.Length {
		messages = append(messages, "пароль должен быть не менее 8 символов")
	}
	if e.HasUpper {
		messages = append(messages, "пароль должен содержать хотя бы одну заглавную букву")
	}
	if e.HasLower {
		messages = append(messages, "пароль должен содержать хотя бы одну строчную букву")
	}
	if e.HasNumber {
		messages = append(messages, "пароль должен содержать хотя бы одну цифру")
	}
	if e.HasSpecial {
		messages = append(messages, "пароль должен содержать хотя бы один специальный символ")
	}

	// Если сообщений нет, значит, ошибок тоже нет.
	if len(messages) == 0 {
		return ""
	}

	return "Пароль не соответствует требованиям: " + strings.Join(messages, ", ")
}

// ValidatePassword проверяет пароль на соответствие критериям безопасности.
// Возвращает true, если пароль валиден, и false вместе со структурой ошибок, если нет.
func ValidatePassword(password string) (bool, PasswordErrors) {
	var (
		errs       PasswordErrors
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	// 1. Проверка длины
	if len([]rune(password)) < 8 {
		errs.Length = true
	}

	// 2. Проверка на наличие разных классов символов
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		// unicode.IsPunct и IsSymbol покрывают большинство спец. символов
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errs.HasUpper = true
	}
	if !hasLower {
		errs.HasLower = true
	}
	if !hasNumber {
		errs.HasNumber = true
	}
	if !hasSpecial {
		errs.HasSpecial = true
	}

	// Если хотя бы одно поле в errs == true, пароль невалиден
	isValid := !errs.Length && !errs.HasUpper && !errs.HasLower && !errs.HasNumber && !errs.HasSpecial

	return isValid, errs
}