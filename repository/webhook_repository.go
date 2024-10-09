package repository

import (
	"log"
	"yayawallet-webhook/models"
)

type WebhookRepository interface {
	Save(payload models.WebhookPayload) error
}

type webhookRepo struct{}

func NewWebhookRepository() WebhookRepository {
	return &webhookRepo{}
}

func (r *webhookRepo) Save(payload models.WebhookPayload) error {
	// Implement saving logic here
	// For now, just log the payload
	log.Printf("Received webhook payload: %+v\n", payload)
	return nil
}
