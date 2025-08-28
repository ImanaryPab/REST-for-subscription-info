package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ServiceName string    `gorm:"not null" json:"service_name"`
	Price       int       `gorm:"not null" json:"price"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     *time.Time `gorm:"null" json:"end_date,omitempty"`
	gorm.Model
}

func (subscription *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	subscription.ID = uuid.New()
	return
}