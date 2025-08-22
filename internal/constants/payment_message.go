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
	InvoiceDeleted = model.Message{
		"en": "Invoice deleted successfully",
		"id": "Tagihan berhasil dihapus",
	}
	OrderNotFound = model.Message{
		"en": "Order not found",
		"id": "Pesanan tidak ditemukan",
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
	InvoiceAlreadyExists = model.Message{
		"en": "Invoice already exists for this order",
		"id": "Tagihan sudah ada untuk pesanan ini",
	}
	UnauthorizedAccess = model.Message{
		"en": "Unauthorized access",
		"id": "Akses tidak sah",
	}
)
