package service

//
//import (
//	"aggregationSubscriptions/internal/models"
//	"aggregationSubscriptions/internal/repository"
//)
//
//type Service interface {
//	GetAllSubscriptions() ([]*models.Subscription, error)
//	GetSubscriptionByID(id string) (*models.Subscription, error)
//	CreateNewSubscription(sub models.Subscription) error
//	UpdateSubscriptionByID(id string, data models.Subscription) (*models.Subscription, error)
//	DeleteSubscriptionByID(id string) error
//	//GetCountSubscriptionsPrice(query) ([]*models.Subscription, error)
//}
//type service struct {
//	repo repository.Repository
//}
//
//func NewService(repo repository.Repository) Service {
//	return &service{repo: repo}
//}
