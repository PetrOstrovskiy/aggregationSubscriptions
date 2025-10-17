package handler

import (
	"aggregationSubscriptions/internal/database"
	"aggregationSubscriptions/internal/models"
	"aggregationSubscriptions/internal/service"
	"aggregationSubscriptions/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSubscriptions(c *gin.Context) {
	var subs []models.Subscription
	db := database.GetDB()

	if err := db.Find(&subs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем в DTO
	var dtoList []models.SubscriptionDTO
	for _, s := range subs {
		dtoList = append(dtoList, models.ToSubscriptionDTO(s))
	}

	c.JSON(http.StatusOK, gin.H{"data": dtoList})
}

func (h *Handler) GetSubscription(c *gin.Context) {
	var sub models.Subscription
	db := database.GetDB()
	id := c.Param("id")

	if err := db.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Подписка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": models.ToSubscriptionDTO(sub)})
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var dto models.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Ошибка записи данных"})
		return
	}

	sub, err := models.ToSubscription(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Ошибка формата данных"})
		return
	}

	if err := utils.ValidateSubscriptionInput(sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Неверные данные подписки"})
		return
	}

	db := database.GetDB()
	if err := db.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось создать запись"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": models.ToSubscriptionDTO(*sub)})
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var existing models.Subscription
	db := database.GetDB()

	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Подписка не найдена"})
		return
	}

	var dto models.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Ошибка записи данных"})
		return
	}

	sub, err := models.ToSubscription(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Ошибка формата данных"})
		return
	}

	// Валидация
	if err := utils.ValidateSubscriptionInput(sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Неверные данные подписки"})
		return
	}

	// Обновляем существующую запись
	existing.ServiceName = sub.ServiceName
	existing.Price = sub.Price
	existing.UserID = sub.UserID
	existing.StartDate = sub.StartDate
	existing.EndDate = sub.EndDate

	if err := db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось сохранить запись"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": models.ToSubscriptionDTO(existing)})
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	if err := db.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *Handler) GetSubscriptionsPrice(c *gin.Context) {
	const monthLayout = "01-2006"
	db := database.GetDB()

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	// Парсинг дат
	start, err := time.Parse(monthLayout, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат start_date, ожидается MM-YYYY"})
		return
	}

	end, err := time.Parse(monthLayout, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат end_date, ожидается MM-YYYY"})
		return
	}

	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date не может быть раньше start_date"})
		return
	}

	query := db.Model(&models.Subscription{}).
		Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", end, start)

	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный user_id"})
			return
		}
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var subs []models.Subscription
	if err := query.Find(&subs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить подписки"})
		return
	}

	// Считаем цену за каждый месяц
	var total int64
	for _, s := range subs {
		endDate := s.EndDate
		if endDate == nil {
			tmp := time.Now()
			endDate = &tmp
		}
		months := utils.MonthsBetween(s.StartDate, *endDate)
		total += int64(s.Price * months)
	}

	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
