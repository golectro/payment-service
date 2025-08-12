package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID       uuid.UUID      `gorm:"type:char(36);index" json:"order_id"`
	UserID        uuid.UUID      `gorm:"type:char(36);index" json:"user_id"`
	XenditID      string         `gorm:"index" json:"xendit_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	PaymentMethod string         `gorm:"size:255" json:"payment_method"`
	PayerEmail    string         `gorm:"size:255" json:"payer_email"`
	Description   string         `gorm:"size:500" json:"description"`
	InvoiceURL    string         `gorm:"size:1000" json:"invoice_url"`
	Status        string         `gorm:"size:50;index" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Invoice) TableName() string {
	return "invoices"
}
