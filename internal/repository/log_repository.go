package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(client *mongo.Client) *LogRepository {
	return &LogRepository{
		collection: client.Database("myapp").Collection("logs"),
	}
}

type LogEntry struct {
	Timestamp time.Time `bson:"timestamp"`
	Level     string    `bson:"level"`
	Message   string    `bson:"message"`
	UserID    string    `bson:"user_id,omitempty"`
	Path      string    `bson:"path"`
}

func (r *LogRepository) Save(ctx context.Context, entry LogEntry) error {
	entry.Timestamp = time.Now()
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}
