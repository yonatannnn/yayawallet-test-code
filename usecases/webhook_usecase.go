package usecases

import (
	"os"
	"time"
	"yayawallet-webhook/models"
)

type WebhookUseCase struct {
	service models.WebhookService
}

func NewWebhookUseCase(service models.WebhookService) *WebhookUseCase {
	return &WebhookUseCase{service: service}
}

var secretKey = os.Getenv("SECRET_KEY")

func (uc *WebhookUseCase) ProcessWebhook(payload models.WebhookPayload, receivedSignature string) (bool, error) {
	if uc.service.VerifySignature(payload, receivedSignature, secretKey) {
		return false, nil
	}

	currentTime := time.Now().Unix()
	if currentTime-payload.Timestamp > 300 {
		return false, nil
	}

	err := uc.service.Save(payload)
	if err != nil {
		return false, err
	}

	return true, nil
}
