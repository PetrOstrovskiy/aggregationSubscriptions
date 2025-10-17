package models

import (
	"fmt"
	"time"
)

type Subscription struct {
	ID          string     `json:"id" gorm:"type:uuid;primaryKey"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionDTO struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

// Конвертация DTO → модель
func ToSubscription(dto SubscriptionDTO) (*Subscription, error) {
	const monthLayout = "01-2006"

	start, err := time.Parse(monthLayout, dto.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date должен быть MM-YYYY")
	}

	var end *time.Time
	if dto.EndDate != nil {
		t, err := time.Parse(monthLayout, *dto.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date должен быть MM-YYYY")
		}
		end = &t
		if end.Before(start) {
			return nil, fmt.Errorf("end_date не может быть раньше start_date")
		}
	}

	return &Subscription{
		ID:          dto.ID,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

// Конвертация модель → DTO
func ToSubscriptionDTO(sub Subscription) SubscriptionDTO {
	const monthLayout = "01-2006"

	dto := SubscriptionDTO{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format(monthLayout),
	}
	if sub.EndDate != nil {
		endStr := sub.EndDate.Format(monthLayout)
		dto.EndDate = &endStr
	}
	return dto
}
