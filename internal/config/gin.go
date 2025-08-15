package config

import (
	"golectro-payment/internal/delivery/http/middleware"
	"golectro-payment/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewGin(viper *viper.Viper, logger *logrus.Logger, mongoDB *mongo.Database, redis *redis.Client) *gin.Engine {
	if viper.GetString("WEB_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	logUC := usecase.NewLogUsecase(mongoDB)

	app := gin.Default()

	app.Use(
		gin.Recovery(),
		SetupCORS(viper),
		middleware.RequestIDMiddleware(),
		middleware.LoggingMiddleware(logger, logUC),
		middleware.NewRateLimiter(viper, redis),
	)

	return app
}
