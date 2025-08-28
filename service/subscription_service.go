package service

import (
	"log"
	"subscription-service/models"
	"subscription-service/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(subscription *models.Subscription) error
	GetSubscription(id uuid.UUID) (*models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id uuid.UUID) error
	ListSubscriptions() ([]models.Subscription, error)
	CalculateTotalCost(userID *uuid.UUID, serviceName *string, startMonthYear, endMonthYear string) (int, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(subscription *models.Subscription) error {
	log.Printf("Service: Creating subscription for user %s", subscription.UserID)
	return s.repo.Create(subscription)
}

func (s *subscriptionService) GetSubscription(id uuid.UUID) (*models.Subscription, error) {
	log.Printf("Service: Getting subscription %s", id)
	return s.repo.GetByID(id)
}

func (s *subscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	log.Printf("Service: Updating subscription %s", subscription.ID)
	return s.repo.Update(subscription)
}

func (s *subscriptionService) DeleteSubscription(id uuid.UUID) error {
	log.Printf("Service: Deleting subscription %s", id)
	return s.repo.Delete(id)
}

func (s *subscriptionService) ListSubscriptions() ([]models.Subscription, error) {
	log.Println("Service: Listing all subscriptions")
	return s.repo.List()
}

func (s *subscriptionService) CalculateTotalCost(userID *uuid.UUID, serviceName *string, startMonthYear, endMonthYear string) (int, error) {
	log.Printf("Service: Calculating total cost from %s to %s", startMonthYear, endMonthYear)
	
	startDate, err := parseMonthYear(startMonthYear)
	if err != nil {
		return 0, err
	}
	
	endDate, err := parseMonthYear(endMonthYear)
	if err != nil {
		return 0, err
	}
	
	return s.repo.CalculateTotalCost(userID, serviceName, startDate, endDate)
}

func parseMonthYear(monthYear string) (time.Time, error) {
	return time.Parse("01-2006", monthYear)
}