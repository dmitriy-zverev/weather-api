package handlers

import "github.com/redis/go-redis/v9"

type WeatherParams struct {
	City string `json:"city"`
}

type ApiConfig struct {
	RedisClient redis.Client
}
