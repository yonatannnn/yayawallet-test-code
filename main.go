package main

import (
	"yayawallet-webhook/controller"
	"yayawallet-webhook/repository"
	"yayawallet-webhook/services"
	"yayawallet-webhook/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set up repository
	repo := repository.NewWebhookRepository()
	services := services.NewWebhookService(repo)
	uc := usecases.NewWebhookUseCase(services)
	handler := controller.NewWebhookHandler(uc)

	// Register the webhook endpoint
	r.POST("/webhook", handler.HandleWebhook)

	r.Run(":3000") // Run on port 8080
}
