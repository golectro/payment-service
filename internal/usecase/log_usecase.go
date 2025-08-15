package usecase

import (
	"context"
	"golectro-payment/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogUseCase struct {
	Collection *mongo.Collection
}

func NewLogUsecase(mongoDB *mongo.Database) *LogUseCase {
	return &LogUseCase{
		Collection: mongoDB.Collection("logs"),
	}
}

func (l *LogUseCase) LogActivity(ctx context.Context, level, requestID, message, userID, endpoint string, statusCode int, errMsg string) error {
	logEntry := model.ActivityLog{
		UserID:     userID,
		Level:      level,
		Message:    message,
		Endpoint:   endpoint,
		StatusCode: statusCode,
		Error:      errMsg,
		RequestID:  requestID,
		Timestamp:  time.Now(),
	}
	_, err := l.Collection.InsertOne(ctx, logEntry)
	return err
}
