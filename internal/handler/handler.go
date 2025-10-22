package handler

import (
	"aggregationSubscriptions/internal/models"
	"aggregationSubscriptions/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// GetSubscriptions godoc
// @Summary      Получить все подписки
// @Description  Возвращает список всех подписок пользователей
// @Tags         subscriptions
// @Produce      json
// @Success      200  {object}  map[string][]models.SubscriptionDTO
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
func (h *Handler) GetSubscriptions(c *gin.Context) {
	subs, err := h.service.GetAllSubscriptions()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось найти записи подписок"})
		return
	}

	slog.Info("Записи были успешно получены")
	c.JSON(http.StatusOK, gin.H{"data": subs})
}

// GetSubscription godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает информацию о конкретной подписке
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки"
// @Success      200  {object}  map[string]models.SubscriptionDTO
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	sub, err := h.service.GetSubscriptionByID(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Запись была успешно получена")
	c.JSON(http.StatusOK, gin.H{"data": sub})
}

// CreateSubscription godoc
// @Summary      Создать новую подписку
// @Description  Добавляет новую подписку в систему
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.SubscriptionDTO  true  "Данные подписки"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscription [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var dto models.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("Ошибка записи данных", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка записи данных"})
		return
	}

	if err := h.service.CreateNewSubscription(dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Запись была успешно создана")
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

// UpdateSubscription godoc
// @Summary      Обновить подписку
// @Description  Изменяет данные существующей подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID подписки"
// @Param        subscription  body  models.SubscriptionDTO  true  "Обновленные данные подписки"
// @Success      200  {object}  map[string]models.SubscriptionDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscription/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var dto models.SubscriptionDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("Ошибка записи данных", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка записи данных"})
		return
	}

	updated, err := h.service.UpdateSubscription(id, dto)
	if err != nil {
		slog.Error("Не удалось сохранить запись", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить запись"})
		return
	}

	slog.Info("Запись была успешно изменена")
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteSubscription godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscription/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	err := h.service.DeleteSubscription(c.Param("id"))

	if err != nil {
		slog.Error("Не удалось удалить запись", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить запись"})
		return
	}

	slog.Info("Запись была успешно удалена")
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

// GetSubscriptionsPrice godoc
// @Summary      Получить общую стоимость подписок
// @Description  Возвращает итоговую стоимость всех подписок по фильтрам
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "ID пользователя"
// @Param        service_name  query     string  false  "Название подписки"
// @Param        start_date    query     string  true   "Дата начала периода"
// @Param        end_date      query     string  true   "Дата конца периода"
// @Success      200  {object}  map[string]int64
// @Failure      400  {object}  map[string]string
// @Router       /subscriptions/aggregate/total [get]
func (h *Handler) GetSubscriptionsPrice(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	total, err := h.service.GetSubscriptionsPrice(userID, serviceName, startStr, endStr)
	if err != nil {
		slog.Error("Не удалось рассчитать итоговую цену", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Итоговая цена успешно получена")
	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
