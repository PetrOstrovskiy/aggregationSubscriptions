package repository

//
//import (
//	"aggregationSubscriptions/internal/models"
//	"gorm.io/gorm"
//)
//
//type Repository interface {
//	GetAllSubscriptions() ([]*models.Subscription, error)
//	GetSubscriptionByID(id string) (*models.Subscription, error)
//	CreateNewSubscription(sub models.Subscription) error
//	UpdateSubscriptionByID(id string, data models.Subscription) (*models.Subscription, error)
//	DeleteSubscriptionByID(id string) error
//	//GetCountSubscriptionsPrice(query) ([]*models.Subscription, error)
//}
//
//type repository struct {
//	db *gorm.DB
//}
//
//func NewRepository(db *gorm.DB) Repository {
//	return &repository{db: db}
//}
