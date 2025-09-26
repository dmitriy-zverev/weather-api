package handlers

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type WeatherForecastType string

const (
	CURRENT      WeatherForecastType = "current"
	N_DAYS       WeatherForecastType = "n_days"
	TODAY_HOURLY WeatherForecastType = "today_hourly"
)

type WeatherParams struct {
	City         string              `json:"city"`
	ForecastType WeatherForecastType `json:"forecast_type"`
	Days         int                 `json:"days"`
}

type ApiConfig struct {
	RedisClient redis.Client
	Limiter     *rate.Limiter
}
