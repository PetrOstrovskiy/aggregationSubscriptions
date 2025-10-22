package service

import (
	"aggregationSubscriptions/internal/models"
	"aggregationSubscriptions/internal/repository"
	"aggregationSubscriptions/internal/utils"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Service interface {
	GetAllSubscriptions() ([]models.SubscriptionDTO, error)
	GetSubscriptionByID(id string) (*models.SubscriptionDTO, error)
	CreateNewSubscription(dto models.SubscriptionDTO) error
	UpdateSubscription(id string, dto models.SubscriptionDTO) (*models.SubscriptionDTO, error)
	DeleteSubscription(id string) error
	GetSubscriptionsPrice(userID, serviceName, startStr, endStr string) (int64, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllSubscriptions() ([]models.SubscriptionDTO, error) {
	sub, err := s.repo.GetAllSubscriptions()

	if err != nil {
		slog.Error("Не удалось найти записи подписок", "error", err)
		return nil, err
	}

	// Конвертируем в DTO
	var dtoList []models.SubscriptionDTO
	for _, s := range sub {
		dtoList = append(dtoList, models.ToSubscriptionDTO(*s))
	}
	return dtoList, nil
}

func (s *service) GetSubscriptionByID(id string) (*models.SubscriptionDTO, error) {
	sub, err := s.repo.GetSubscriptionByID(id)
	if err != nil {
		slog.Error("Не удалось найти подписку", "error", err)
		return nil, err
	}
	dto := models.ToSubscriptionDTO(*sub)
	return &dto, nil
}

func (s *service) CreateNewSubscription(dto models.SubscriptionDTO) error {
	sub, err := models.ToSubscription(dto)
	if err != nil {
		slog.Error("Ошибка с форматом данных даты", "error", err)
		return err
	}

	sub.ID = uuid.New().String()

	if err := utils.ValidateSubscription(sub); err != nil {
		slog.Error("Не удалось создать запись", "error", err)
		return err
	}
	return s.repo.CreateNewSubscription(sub)
}

func (s *service) UpdateSubscription(id string, dto models.SubscriptionDTO) (*models.SubscriptionDTO, error) {
	sub, err := models.ToSubscription(dto)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateSubscription(sub); err != nil {
		return nil, err
	}

	updatedSub, err := s.repo.UpdateSubscriptionByID(id, sub)
	if err != nil {
		return nil, err
	}

	dtoResponse := models.ToSubscriptionDTO(*updatedSub)
	return &dtoResponse, nil

}

func (s *service) DeleteSubscription(id string) error {
	return s.repo.DeleteSubscriptionByID(id)
}

func (s *service) GetSubscriptionsPrice(userID, serviceName, startStr, endStr string) (int64, error) {
	const monthLayout = "01-2006"

	// Парсинг дат
	start, err := time.Parse(monthLayout, startStr)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse(monthLayout, endStr)
	if err != nil {
		return 0, err
	}

	if end.Before(start) {
		return 0, errors.New("end_date не может быть раньше start_date")
	}

	subs, err := s.repo.GetCountSubscriptionsPrice(userID, serviceName, start, end)
	if err != nil {
		return 0, err
	}

	// Считаем итоговую цену
	var total int64
	for _, sub := range subs {
		actualEnd := sub.EndDate
		if actualEnd == nil || actualEnd.After(end) {
			actualEnd = &end
		}

		months := utils.MonthsBetween(sub.StartDate, *actualEnd)
		total += int64(sub.Price * months)
	}

	return total, nil

}
