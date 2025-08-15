package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Auth struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username" validate:"required,min=3,max=50"`
	Email    string         `json:"email" validate:"required,email"`
	Roles    datatypes.JSON `json:"roles" validate:"required,dive,required"`
}
