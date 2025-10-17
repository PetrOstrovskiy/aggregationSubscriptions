package database

import (
	"aggregationSubscriptions/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

var db *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Не удалось подключиться к базе данных",
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
		slog.Error("Попытка обращения к базе до инициализации!")
		os.Exit(1)
	}
	return db
}
