package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupCORS(viper *viper.Viper) gin.HandlerFunc {
	allowOrigins := viper.GetStringSlice("CORS_ALLOW_ORIGINS")
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers", "X-CSRF-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Requested-With", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	})
}
