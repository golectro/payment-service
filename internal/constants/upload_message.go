package constants

import "golectro-payment/internal/model"

var (
	InvalidFileType = model.Message{
		"en": "Invalid file type",
		"id": "Tipe file tidak valid",
	}
	InvalidResetPosition = model.Message{
		"en": "Failed to reset file position",
		"id": "Gagal mengatur ulang posisi file",
	}
	InvalidReadFile = model.Message{
		"en": "Failed to read file",
		"id": "Gagal membaca file",
	}
	InvalidOpenFile = model.Message{
		"en": "Failed to open file",
		"id": "Gagal membuka file",
	}
	FileSizeExceeded = model.Message{
		"en": "File size exceeded",
		"id": "Ukuran file melebihi batas",
	}
	FileNotFound = model.Message{
		"en": "File not found",
		"id": "File tidak ditemukan",
	}
	FailedUploadObject = model.Message{
		"en": "Failed to upload avatar",
		"id": "Gagal mengunggah avatar",
	}
)
