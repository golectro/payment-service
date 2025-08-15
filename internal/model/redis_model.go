package model

type RedisRequest struct {
	Value string `json:"value" validate:"required"`
	TTL   int    `json:"ttl" validate:"required,min=0"`
}
