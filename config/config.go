// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// GetConfig загружает конфигурацию из .env и возвращает полностью инициализированный объект *pgx.ConnConfig.
func GetConfig() *pgx.ConnConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 1. Собираем DSN-строку из отдельных переменных окружения.
	// Это необходимо, чтобы передать ее в официальный конструктор pgx.ParseConfig.
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", ""),
		getEnv("DB_SSLMODE", "disable"),
	)

	// 2. ИСПОЛЬЗУЕМ ОБЯЗАТЕЛЬНЫЙ КОНСТРУКТОР ParseConfig.
	// Он создает и правильно инициализирует структуру.
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DSN: %v", err)
	}

	return config
}

// Вспомогательная функция остается без изменений.
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			return defaultValue
		}
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}