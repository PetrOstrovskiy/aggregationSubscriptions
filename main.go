package main

import (
	"aggregationSubscriptions/internal/database"
	"aggregationSubscriptions/internal/handler"
	"aggregationSubscriptions/internal/repository"
	"aggregationSubscriptions/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	database.Connect()
	database.Migrate()

	db := database.GetDB()
	subRepository := repository.NewRepository(db)
	subService := service.NewService(subRepository)
	subHandler := handler.NewHandler(subService)

	router := gin.Default()

	router.GET("/subscriptions", subHandler.GetSubscriptions)
	router.GET("/subscriptions/:id", subHandler.GetSubscription)
	router.POST("/subscription", subHandler.CreateSubscription)
	router.PUT("/subscription/:id", subHandler.UpdateSubscription)
	router.DELETE("/subscription/:id", subHandler.DeleteSubscription)
	router.GET("/subscriptions/aggregate/total", subHandler.GetSubscriptionsPrice)

	slog.Info("Сервер запущен на http://localhost:8080")
	router.Run(":8080")
}
