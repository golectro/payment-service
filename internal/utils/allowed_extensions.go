package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

var allowedExtensions = []string{".pdf", ".jpg", ".jpeg", ".png"}

func ValidateFileExtension(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if slices.Contains(allowedExtensions, ext) {
		return nil
	}
	return errors.New("file type not allowed: " + ext)
}
