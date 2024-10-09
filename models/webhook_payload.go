package models

type WebhookPayload struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	CreatedAt   int64  `json:"created_at_time"`
	Timestamp   int64  `json:"timestamp"`
	Cause       string `json:"cause"`
	FullName    string `json:"full_name"`
	AccountName string `json:"account_name"`
	InvoiceURL  string `json:"invoice_url"`
}
