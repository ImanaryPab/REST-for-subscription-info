package main

import (
	"log"
	"subscription-service/config"
	"subscription-service/controllers"
	"subscription-service/database"
	"subscription-service/repository"
	"subscription-service/service"

	"github.com/gin-gonic/gin"
	_ "subscription-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	cfg := config.LoadConfig()

	db := database.InitDB(cfg)
	
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", subscriptionController.CreateSubscription)
		api.GET("/subscriptions", subscriptionController.ListSubscriptions)
		api.GET("/subscriptions/:id", subscriptionController.GetSubscription)
		api.PUT("/subscriptions/:id", subscriptionController.UpdateSubscription)
		api.DELETE("/subscriptions/:id", subscriptionController.DeleteSubscription)
		api.GET("/cost", subscriptionController.CalculateTotalCost)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}