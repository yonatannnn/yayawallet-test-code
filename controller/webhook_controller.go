package controller

import (
	"net/http"
	"yayawallet-webhook/models"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	usecase models.WebhookUseCase
}

func NewWebhookHandler(uc models.WebhookUseCase) *WebhookHandler {
	return &WebhookHandler{usecase: uc}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	var payload models.WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	signature := c.GetHeader("YAYA-SIGNATURE")
	success, err := h.usecase.ProcessWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if !success {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature or request is too old"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
