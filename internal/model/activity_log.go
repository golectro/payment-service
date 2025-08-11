package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	Level      string             `bson:"level"`
	Message    string             `bson:"message"`
	Endpoint   string             `bson:"endpoint"`
	StatusCode int                `bson:"status_code"`
	Error      string             `bson:"error,omitempty"`
	RequestID  string             `bson:"request_id,omitempty"`
	Timestamp  time.Time          `bson:"timestamp"`
}
