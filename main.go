package main

import (
	_ "aggregationSubscriptions/docs"
	"aggregationSubscriptions/internal/database"
	"aggregationSubscriptions/internal/handler"
	"aggregationSubscriptions/internal/repository"
	"aggregationSubscriptions/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log/slog"
	"os"
)

// @title Aggregation Subscriptions API
// @version 1.0
// @description API для управления подписками пользователей.
// @host localhost:8080
// @BasePath /
func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	database.Connect()
	database.Migrate()

	db := database.GetDB()
	subRepository := repository.NewRepository(db)
	subService := service.NewService(subRepository)
	subHandler := handler.NewHandler(subService)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/subscriptions", subHandler.GetSubscriptions)
	router.GET("/subscription/:id", subHandler.GetSubscription)
	router.POST("/subscription", subHandler.CreateSubscription)
	router.PUT("/subscription/:id", subHandler.UpdateSubscription)
	router.DELETE("/subscription/:id", subHandler.DeleteSubscription)
	router.GET("/subscriptions/aggregate/total", subHandler.GetSubscriptionsPrice)

	slog.Info("Сервер запущен на http://localhost:8080")
	router.Run(":8080")
}
