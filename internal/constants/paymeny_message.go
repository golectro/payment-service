package constants

import "golectro-payment/internal/model"

var (
	InvoiceCreated = model.Message{
		"en": "Invoice created successfully",
		"id": "Tagihan berhasil dibuat",
	}
	InvoiceRetrieved = model.Message{
		"en": "Invoice retrieved successfully",
		"id": "Tagihan berhasil diambil",
	}
)

var (
	InvalidPaymentMethod = model.Message{
		"en": "Invalid payment method",
		"id": "Metode pembayaran tidak valid",
	}
	FailedToCreateInvoice = model.Message{
		"en": "Failed to create invoice",
		"id": "Gagal membuat tagihan",
	}
	UnauthorizedAccess = model.Message{
		"en": "Unauthorized access",
		"id": "Akses tidak sah",
	}
)
