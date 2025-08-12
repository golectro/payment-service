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
