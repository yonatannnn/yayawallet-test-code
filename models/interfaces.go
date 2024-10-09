package models

// WebhookPayload represents the payload of a webhook request.
type WebhookRepository interface {
	Save(payload WebhookPayload) error
}

type WebhookService interface {
	VerifySignature(payload WebhookPayload, receivedSignature string, secretKey string) bool
	CreateSignedPayload(payload WebhookPayload) string
	GenerateHMAC(data, secret string) string
	Save(payload WebhookPayload) error
}

type WebhookUseCase interface {
	ProcessWebhook(payload WebhookPayload, receivedSignature string) (bool, error)
}
