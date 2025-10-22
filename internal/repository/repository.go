package repository

import (
	"aggregationSubscriptions/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	GetAllSubscriptions() ([]*models.Subscription, error)
	GetSubscriptionByID(id string) (*models.Subscription, error)
	CreateNewSubscription(sub *models.Subscription) error
	UpdateSubscriptionByID(id string, data *models.Subscription) (*models.Subscription, error)
	DeleteSubscriptionByID(id string) error
	GetCountSubscriptionsPrice(userID string, serviceName string, start time.Time, end time.Time) ([]*models.Subscription, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllSubscriptions() ([]*models.Subscription, error) {
	var subs []*models.Subscription

	err := r.db.Find(&subs).Error
	return subs, err
}

func (r *repository) GetSubscriptionByID(id string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.First(&sub, "id = ?", id).Error
	return &sub, err
}

func (r *repository) CreateNewSubscription(sub *models.Subscription) error {
	err := r.db.Create(sub).Error
	return err
}

func (r *repository) UpdateSubscriptionByID(id string, data *models.Subscription) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := r.db.First(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}
	subscription.ServiceName = data.ServiceName
	subscription.Price = data.Price
	subscription.UserID = data.UserID
	subscription.StartDate = data.StartDate
	subscription.EndDate = data.EndDate

	if err := r.db.Save(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *repository) DeleteSubscriptionByID(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Subscription{}).Error
	return err
}

func (r *repository) GetCountSubscriptionsPrice(userID string, serviceName string, start time.Time, end time.Time) ([]*models.Subscription, error) {
	var subs []*models.Subscription
	query := r.db.Model(&models.Subscription{}).Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", end, start)

	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			return nil, err
		}
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}
