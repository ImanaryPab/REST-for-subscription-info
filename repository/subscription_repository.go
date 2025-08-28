package repository

import (
	"log"
	"subscription-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uuid.UUID) (*models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id uuid.UUID) error
	List() ([]models.Subscription, error)
	CalculateTotalCost(userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	log.Printf("Creating subscription for user %s to service %s", subscription.UserID, subscription.ServiceName)
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	return &subscription, err
}

func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	log.Printf("Updating subscription %s", subscription.ID)
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	log.Printf("Deleting subscription %s", id)
	return r.db.Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepository) List() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) CalculateTotalCost(userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int, error) {
	log.Printf("Calculating total cost for period %s to %s", startDate.Format("01-2006"), endDate.Format("01-2006"))
	
	query := r.db.Model(&models.Subscription{}).
		Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endDate, startDate)
	
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	
	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}
	
	var totalCost int
	err := query.Select("COALESCE(SUM(price), 0)").Row().Scan(&totalCost)
	
	return totalCost, err
}