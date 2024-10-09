package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"yayawallet-webhook/models"
)

type WebhookService struct {
	repo models.WebhookRepository
}

func NewWebhookService(repo models.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

func (ws *WebhookService) VerifySignature(payload models.WebhookPayload, receivedSignature string, secretKey string) bool {
	signedPayload := ws.CreateSignedPayload(payload)
	expectedSignature := ws.GenerateHMAC(signedPayload, secretKey)
	return hmac.Equal([]byte(receivedSignature), []byte(expectedSignature))
}

func (ws *WebhookService) CreateSignedPayload(payload models.WebhookPayload) string {
	return payload.ID + strconv.Itoa(payload.Amount) + payload.Currency +
		strconv.FormatInt(payload.CreatedAt, 10) + // Convert int64 to string
		strconv.FormatInt(payload.Timestamp, 10) + // Convert int64 to string
		payload.Cause + payload.FullName + payload.AccountName +
		payload.InvoiceURL
}

func (ws *WebhookService) GenerateHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (ws *WebhookService) Save(payload models.WebhookPayload) error {
	return ws.repo.Save(payload)
}
