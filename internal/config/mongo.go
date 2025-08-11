package config

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(viper *viper.Viper, log *logrus.Logger) *mongo.Database {
	uri := viper.GetString("MONGO_URI")
	dbName := viper.GetString("MONGO_DB")
	if uri == "" || dbName == "" {
		log.Fatal("MONGO_URI and MONGO_DB are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB not responding: %v", err)
	}

	log.Infof("MongoDB connected to %s", uri)
	return client.Database(dbName)
}
