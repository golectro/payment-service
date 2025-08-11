package config

import (
	"golectro-payment/internal/delivery/http/middleware"
	"golectro-payment/internal/delivery/http/route"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/vault/api"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	Mongo      *mongo.Database
	App        *gin.Engine
	Redis      *redis.Client
	Log        *logrus.Logger
	Validate   *validator.Validate
	Viper      *viper.Viper
	GRPCClient *grpc.ClientConn
	Vault      *api.Client
}

func Bootstrap(config *BootstrapConfig) {
	authMiddleware := middleware.NewAuth(config.Viper)

	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		Viper:          config.Viper,
	}
	routeConfig.Setup()
}
