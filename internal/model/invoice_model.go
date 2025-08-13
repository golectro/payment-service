package model

type CreateInvoiceRequest struct {
	OrderID     string  `json:"order_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description" validate:"required"`
}

type CreateInvoiceResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	XenditID   string  `json:"xendit_id"`
	InvoiceURL string  `json:"invoice_url"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

type InvoiceResponse struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	XenditID    string  `json:"xendit_id"`
	InvoiceURL  string  `json:"invoice_url"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	PayerEmail  string  `json:"payer_email"`
	Description string  `json:"description"`
}

type XenditCallbackData struct {
	ID             string  `json:"id" validate:"required"`
	ExternalID     string  `json:"external_id" validate:"required"`
	Amount         float64 `json:"amount" validate:"required"`
	Status         string  `json:"status" validate:"required"`
	PayerEmail     string  `json:"payer_email" validate:"required,email"`
	Description    string  `json:"description" validate:"required"`
	PaymentMethod  string  `json:"payment_method" validate:"required"`
	PaymentChannel string  `json:"payment_channel" validate:"required"`
}
