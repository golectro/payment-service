package repository

import (
	"golectro-payment/internal/entity"

	"github.com/sirupsen/logrus"
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
