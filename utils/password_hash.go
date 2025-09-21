package utils

import "golang.org/x/crypto/bcrypt"

// GetPasswordHash создает безопасный bcrypt хэш из строки пароля.
// Эта функция будет использоваться при регистрации пользователя или смене пароля.
func GetPasswordHash(password string) (string, error) {
	// bcrypt.GenerateFromPassword принимает пароль в виде среза байтов.
	// Второй аргумент - это "стоимость" (cost) хэширования.
	// Чем она выше, тем дольше вычисляется хэш и тем сложнее его
	// взломать методом перебора (брутфорсом).
	// bcrypt.DefaultCost - это хороший, сбалансированный выбор (сейчас это 10).
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Если при хэшировании произошла ошибка, возвращаем ее.
		return "", err
	}

	// Конвертируем срез байтов с хэшем в строку и возвращаем.
	return string(hashedBytes), nil
}

// CheckPasswordHash сравнивает предоставленный пароль с хэшем из базы данных.
// Эта функция будет использоваться при логине пользователя.
func CheckPasswordHash(password, hash string) bool {
	// bcrypt.CompareHashAndPassword выполняет ту же самую "медленную" операцию
	// сравнения, что делает его устойчивым к атакам по времени (timing attacks).
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// Если err == nil, значит пароль и хэш совпали.
	// В противном случае, err будет bcrypt.ErrMismatchedHashAndPassword.
	return err == nil
}