package database

import (
	"aggregationSubscriptions/internal/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

var db *gorm.DB

func Connect() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("Не удалось загрузить .env файл, используются системные переменные")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	// Проверка на пустоту переменных для подключения к ДБ
	if host == "" || port == "" || user == "" || name == "" {
		slog.Error("Не заданы переменные окружения для подключения к PostgreSQL")
		os.Exit(1)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Не удалось подключиться к БД",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Успешное подключение к БД")

}

func Migrate() {
	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		slog.Error("Ошибка миграции", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		slog.Error("Попытка обращения к базе до инициализации")
		os.Exit(1)
	}
	return db
}
