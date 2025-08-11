package utils

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func GetHTTPStatusCode(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
