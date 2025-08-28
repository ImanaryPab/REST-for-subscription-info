package controllers

import (
	"log"
	"net/http"
	"subscription-service/models"
	"subscription-service/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionController struct {
	service service.SubscriptionService
}

func NewSubscriptionController(service service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{service: service}
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (c *SubscriptionController) CreateSubscription(ctx *gin.Context) {
	log.Println("Controller: Creating subscription")
	
	var req CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Printf("Error parsing start date: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			log.Printf("Error parsing end date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use MM-YYYY"})
			return
		}
		endDate = &parsedEndDate
	}

	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := c.service.CreateSubscription(subscription); err != nil {
		log.Printf("Error creating subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	log.Printf("Subscription created successfully: %s", subscription.ID)
	ctx.JSON(http.StatusCreated, subscription)
}

// GetSubscription godoc
// @Summary Get a subscription by ID
// @Description Get subscription details by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (c *SubscriptionController) GetSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	log.Printf("Controller: Getting subscription %s", id)
	subscription, err := c.service.GetSubscription(id)
	if err != nil {
		log.Printf("Error getting subscription: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update subscription details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body models.Subscription true "Subscription data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (c *SubscriptionController) UpdateSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	log.Printf("Controller: Updating subscription %s", id)
	var subscription models.Subscription
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription.ID = id
	if err := c.service.UpdateSubscription(&subscription); err != nil {
		log.Printf("Error updating subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (c *SubscriptionController) DeleteSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	log.Printf("Controller: Deleting subscription %s", id)
	if err := c.service.DeleteSubscription(id); err != nil {
		log.Printf("Error deleting subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListSubscriptions godoc
// @Summary List all subscriptions
// @Description Get a list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (c *SubscriptionController) ListSubscriptions(ctx *gin.Context) {
	log.Println("Controller: Listing all subscriptions")
	subscriptions, err := c.service.ListSubscriptions()
	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

type CalculateCostRequest struct {
	UserID      *uuid.UUID `form:"user_id"`
	ServiceName *string    `form:"service_name"`
	StartDate   string     `form:"start_date" binding:"required"`
	EndDate     string     `form:"end_date" binding:"required"`
}

// CalculateTotalCost godoc
// @Summary Calculate total cost of subscriptions
// @Description Calculate total cost of subscriptions for a given period with optional filters
// @Tags cost
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Param start_date query string true "Start Date (MM-YYYY)"
// @Param end_date query string true "End Date (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cost [get]
func (c *SubscriptionController) CalculateTotalCost(ctx *gin.Context) {
	log.Println("Controller: Calculating total cost")
	
	var req CalculateCostRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding query: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalCost, err := c.service.CalculateTotalCost(req.UserID, req.ServiceName, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Error calculating cost: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate cost"})
		return
	}

	log.Printf("Total cost calculated: %d", totalCost)
	ctx.JSON(http.StatusOK, gin.H{"total_cost": totalCost})
}