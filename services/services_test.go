package services

import (
	"testing"
	"time"
	"yayawallet-webhook/models"
	"yayawallet-webhook/models/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebhookServiceSuite struct {
	suite.Suite
	repoMock *mocks.WebhookRepository
	ws       *WebhookService
}

func (suite *WebhookServiceSuite) SetupSuite() {
	// Initialize the mock repository and the service before any tests run
	suite.repoMock = new(mocks.WebhookRepository)
	suite.ws = NewWebhookService(suite.repoMock)
}

func (suite *WebhookServiceSuite) TearDownSuite() {
	// Clean up after all tests are run
	suite.repoMock = nil
	suite.ws = nil
}

func (suite *WebhookServiceSuite) TestVerifySignature_ValidSignature() {
	// Setup
	secretKey := "mysecret"
	payload := models.WebhookPayload{
		ID:          "12345",
		Amount:      1000,
		Currency:    "USD",
		CreatedAt:   time.Now().Unix(),
		Timestamp:   time.Now().Unix(),
		Cause:       "Payment",
		FullName:    "John Doe",
		AccountName: "john.doe@example.com",
		InvoiceURL:  "http://example.com/invoice/12345",
	}

	signedPayload := suite.ws.CreateSignedPayload(payload)
	expectedSignature := suite.ws.GenerateHMAC(signedPayload, secretKey)

	// Act
	isValid := suite.ws.VerifySignature(payload, expectedSignature, secretKey)

	// Assert
	assert.True(suite.T(), isValid, "The signature should be valid")
}

func (suite *WebhookServiceSuite) TestVerifySignature_InvalidSignature() {
	// Setup
	secretKey := "mysecret"
	payload := models.WebhookPayload{
		ID:          "12345",
		Amount:      1000,
		Currency:    "USD",
		CreatedAt:   time.Now().Unix(),
		Timestamp:   time.Now().Unix(),
		Cause:       "Payment",
		FullName:    "John Doe",
		AccountName: "john.doe@example.com",
		InvoiceURL:  "http://example.com/invoice/12345",
	}

	// Act
	isValid := suite.ws.VerifySignature(payload, "invalid_signature", secretKey)

	// Assert
	assert.False(suite.T(), isValid, "The signature should be invalid")
}

func (suite *WebhookServiceSuite) TestCreateSignedPayload() {
	// Setup
	payload := models.WebhookPayload{
		ID:          "12345",
		Amount:      1000,
		Currency:    "USD",
		CreatedAt:   1625097600,
		Timestamp:   1625097600,
		Cause:       "Payment",
		FullName:    "John Doe",
		AccountName: "john.doe@example.com",
		InvoiceURL:  "http://example.com/invoice/12345",
	}

	expectedPayload := "123451000USD16250976001625097600PaymentJohn Doejohn.doe@example.comhttp://example.com/invoice/12345"

	// Act
	actualPayload := suite.ws.CreateSignedPayload(payload)

	// Assert
	assert.Equal(suite.T(), expectedPayload, actualPayload, "The signed payload should match the expected value")
}

func (suite *WebhookServiceSuite) TestSave() {
	// Setup
	payload := models.WebhookPayload{
		ID:          "12345",
		Amount:      1000,
		Currency:    "USD",
		CreatedAt:   time.Now().Unix(),
		Timestamp:   time.Now().Unix(),
		Cause:       "Payment",
		FullName:    "John Doe",
		AccountName: "john.doe@example.com",
		InvoiceURL:  "http://example.com/invoice/12345",
	}

	// Expect the Save method to be called
	suite.repoMock.On("Save", payload).Return(nil)

	// Act
	err := suite.ws.Save(payload)

	// Assert
	assert.NoError(suite.T(), err, "Save should not return an error")
	suite.repoMock.AssertExpectations(suite.T())
}

func TestWebhookServiceSuite(t *testing.T) {
	suite.Run(t, new(WebhookServiceSuite))
}
