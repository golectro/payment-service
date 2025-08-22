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

func (r *InvoiceRepository) FindByOrderID(tx *gorm.DB, orderID uuid.UUID, invoice *entity.Invoice) error {
	if err := tx.Where("order_id = ?", orderID).Where("status = ?", "PENDING").First(invoice).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find invoice by order ID")
		return err
	}
	return nil
}

func (r *InvoiceRepository) UpdateInvoice(
	tx *gorm.DB, orderID uuid.UUID, xenditID string, invoice *entity.Invoice,
) error {
	r.Log.Infof("Updating invoice with orderID: %s and xenditID: %s", orderID, xenditID)

	result := tx.Model(&entity.Invoice{}).
		Where("order_id = ? AND xendit_id = ?", orderID, xenditID).
		Omit("id", "created_at").
		Updates(invoice)

	if result.Error != nil {
		r.Log.WithError(result.Error).Error("Failed to update invoice")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.Log.Warn("No invoice updated — possible invalid orderID")
	}
	return nil
}

func (r *InvoiceRepository) UpdateDeleteColumn(tx *gorm.DB, userID uuid.UUID, xenditID string) error {
	r.Log.Infof("Updating delete column for userID: %s and xenditID: %s", userID, xenditID)

	result := tx.Model(&entity.Invoice{}).
		Where("user_id = ? AND xendit_id = ?", userID, xenditID).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		r.Log.WithError(result.Error).Error("Failed to update delete column")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.Log.Warn("No invoice updated — possible invalid userID or xenditID")
	}
	return nil
}

func (r *InvoiceRepository) FindByUserIDAndXenditID(tx *gorm.DB, userID uuid.UUID, xenditID string, invoice *entity.Invoice) error {
	if err := tx.Where("user_id = ? AND xendit_id = ?", userID, xenditID).First(invoice).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find invoice by user ID and Xendit ID")
		return err
	}
	return nil
}

func (r *InvoiceRepository) FindAllExceptDeletedByUserID(tx *gorm.DB, userID uuid.UUID, invoices *[]entity.Invoice) error {
	if err := tx.Where("user_id = ? AND deleted_at IS NULL", userID).Find(invoices).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find all invoices except deleted by user ID")
		return err
	}
	return nil
}
