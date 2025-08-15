package middleware

import (
	"fmt"
	"golectro-payment/internal/constants"
	"golectro-payment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewRateLimiter(viper *viper.Viper, redis *redis.Client) gin.HandlerFunc {
	rateStr := viper.GetString("RATE_LIMIT")
	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		panic(err)
	}

	store, err := redisstore.NewStoreWithOptions(redis, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		panic(err)
	}

	limiterInstance := limiter.New(store, rate)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}

		limiterCtx, err := limiterInstance.Get(ctx, ip)
		if err != nil {
			res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterCtx.Limit))
		ctx.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterCtx.Remaining))
		ctx.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterCtx.Reset))

		if limiterCtx.Reached {
			res := utils.FailedResponse(ctx, http.StatusTooManyRequests, constants.TooManyRequests, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		ctx.Next()
	}
}
