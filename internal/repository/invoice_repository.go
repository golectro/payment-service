package repository

import (
	"golectro-payment/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	Repository[entity.Invoice]
	Log *logrus.Logger
}

func NewInvoiceRepository(log *logrus.Logger) *InvoiceRepository {
	return &InvoiceRepository{
		Log: log,
	}
}

func (r *InvoiceRepository) FindByUserID(tx *gorm.DB, userID uuid.UUID, invoice *entity.Invoice) error {
	if err := tx.Where("user_id = ?", userID).First(invoice).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find invoice by user ID")
		return err
	}
	return nil
}
