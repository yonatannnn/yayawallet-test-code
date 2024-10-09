package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"yayawallet-webhook/models"
	"yayawallet-webhook/models/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebhookHandlerSuite struct {
	suite.Suite
	usecaseMock *mocks.WebhookUseCase
	handler     *WebhookHandler
}

func (suite *WebhookHandlerSuite) SetupTest() {
	// Reset the mock before each test to ensure isolation
	suite.usecaseMock = new(mocks.WebhookUseCase)
	suite.handler = NewWebhookHandler(suite.usecaseMock)

	// Clean up and assert expectations after each test
	suite.T().Cleanup(func() {
		suite.usecaseMock.AssertExpectations(suite.T())
	})
}

func (suite *WebhookHandlerSuite) TestHandleWebhook_Success() {
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

	suite.usecaseMock.On("ProcessWebhook", payload, "valid_signature").Return(true, nil)

	// Prepare the request body
	body := `{"id":"12345","amount":1000,"currency":"USD","created_at_time":1625097600,"timestamp":1625097600,"cause":"Payment","full_name":"John Doe","account_name":"john.doe@example.com","invoice_url":"http://example.com/invoice/12345"}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(body))
	req.Header.Set("YAYA-SIGNATURE", "valid_signature")

	// Create a gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	suite.handler.HandleWebhook(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"status":"success"}`, w.Body.String())
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerSuite) TestHandleWebhook_InvalidJSON() {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	suite.handler.HandleWebhook(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(suite.T(), `{"error":"Invalid JSON"}`, w.Body.String())
}

func (suite *WebhookHandlerSuite) TestHandleWebhook_ProcessWebhookError() {
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

	suite.usecaseMock.On("ProcessWebhook", payload, "valid_signature").Return(false, nil)

	// Prepare the request
	body := `{
	"id": "12345",
	"amount": 1000,
	"currency": "USD",
	"created_at_time": 1625097600,
	"timestamp": 1625097600,
	"cause": "Payment",
	"full_name": "John Doe",
	"account_name": "john.doe@example.com",
	"invoice_url": "http://example.com/invoice/12345"
}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(body))
	req.Header.Set("YAYA-SIGNATURE", "valid_signature")

	// Create a gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	suite.handler.HandleWebhook(c)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.JSONEq(suite.T(), `{"error":"Invalid signature or request is too old"}`, w.Body.String())
}

func (suite *WebhookHandlerSuite) TestHandleWebhook_InternalServerError() {
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

	suite.usecaseMock.On("ProcessWebhook", payload, "valid_signature").Return(false, assert.AnError)

	// Prepare the request
	body := `{"id":"12345","amount":1000,"currency":"USD","created_at_time":1625097600,"timestamp":1625097600,"cause":"Payment","full_name":"John Doe","account_name":"john.doe@example.com","invoice_url":"http://example.com/invoice/12345"}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(body))
	req.Header.Set("YAYA-SIGNATURE", "valid_signature")

	// Create a gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	suite.handler.HandleWebhook(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"error":"Internal Server Error"}`, w.Body.String())
}

func TestWebhookHandlerSuite(t *testing.T) {
	suite.Run(t, new(WebhookHandlerSuite))
}
