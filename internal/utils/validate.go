package utils

import (
	"aggregationSubscriptions/internal/models"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

const monthLayout = "01-2006"

func ValidateSubscription(sub *models.Subscription) error {
	sub.ServiceName = strings.TrimSpace(sub.ServiceName)
	sub.UserID = strings.TrimSpace(sub.UserID)

	if sub.ServiceName == "" {
		return fmt.Errorf("service_name обязателен")
	}

	if sub.Price <= 0 {
		return fmt.Errorf("price должен быть > 0")
	}

	if _, err := uuid.Parse(sub.UserID); err != nil {
		return fmt.Errorf("неверный user_id")
	}

	startStr := sub.StartDate.Format(monthLayout)
	_, err := time.Parse(monthLayout, startStr)
	if err != nil {
		return fmt.Errorf("start_date должен быть в формате MM-YYYY")
	}

	if sub.EndDate != nil {
		endStr := sub.EndDate.Format(monthLayout)
		end, err := time.Parse(monthLayout, endStr)
		if err != nil {
			return fmt.Errorf("end_date должен быть в формате MM-YYYY")
		}
		if end.Before(sub.StartDate) {
			return fmt.Errorf("end_date не может быть раньше start_date")
		}
	}

	return nil
}

func MonthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	return years*12 + months + 1
}
