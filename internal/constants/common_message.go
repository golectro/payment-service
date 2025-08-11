package constants

import "golectro-payment/internal/model"

var (
	WelcomeMessage = model.Message{
		"en": "Welcome to Golectro Payment API",
		"id": "Selamat datang di Golectro Payment API",
	}
	NotFound = model.Message{
		"en": "API not found",
		"id": "API tidak ditemukan",
	}
	FailedDataFromBody = model.Message{
		"en": "Failed to get data from body",
		"id": "Gagal mendapatkan data dari body",
	}
	FailedInputFormat = model.Message{
		"en": "Invalid input format",
		"id": "Format input tidak valid",
	}
	FailedValidationOccurred = model.Message{
		"en": "Validation error occurred",
		"id": "Terjadi kesalahan validasi",
	}
	InvalidToken = model.Message{
		"en": "Invalid token",
		"id": "Token tidak valid",
	}
	InvalidRequestData = model.Message{
		"en": "Invalid request data",
		"id": "Data permintaan tidak valid",
	}
	InternalServerError = model.Message{
		"en": "Internal server error",
		"id": "Kesalahan internal server",
	}
	TooManyRequests = model.Message{
		"en": "Too many requests, please try again later",
		"id": "Terlalu banyak permintaan, silakan coba lagi nanti",
	}
)
