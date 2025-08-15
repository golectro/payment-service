package usecase

import (
	"context"
	"golectro-payment/internal/model"
	"golectro-payment/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisUseCase struct {
	Client   *redis.Client
	Validate *validator.Validate
	Log      *logrus.Logger
}

func NewRedisUsecase(client *redis.Client, log *logrus.Logger, validate *validator.Validate) *RedisUseCase {
	return &RedisUseCase{
		Client:   client,
		Validate: validate,
		Log:      log,
	}
}

func (r *RedisUseCase) ValidateRequest(request *model.RedisRequest) error {
	if err := r.Validate.Struct(request); err != nil {
		r.Log.WithError(err).Error("Invalid input format")
		message := utils.TranslateValidationError(r.Validate, err)
		return utils.WrapMessageAsError(message)
	}
	return nil
}

func (r *RedisUseCase) Set(ctx context.Context, key, value string, ttlSeconds int) error {
	if ttlSeconds <= 0 {
		ttlSeconds = 60
	}
	exp := time.Duration(ttlSeconds) * time.Second
	return r.Client.Set(ctx, "cache:"+key, value, exp).Err()
}

func (r *RedisUseCase) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, "cache:"+key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisUseCase) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, "cache:"+key).Err()
}
