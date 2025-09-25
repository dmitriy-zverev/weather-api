package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dmitriy-zverev/weather-api/internal/handlers"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisURL string
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("error: ", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         config.RedisURL,
		Password:     "",
		DB:           0,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
	})

	apiConfig := &handlers.ApiConfig{
		RedisClient: *rdb,
	}

	mux := setupRoutes(apiConfig)

	if err := startServer(mux, PORT); err != nil {
		log.Fatal("server error:", err)
	}
}

func setupRoutes(cfg *handlers.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+BASE_ROUTE, handlers.BaseHandler)
	mux.HandleFunc("GET "+WEATHER_ROUTE, cfg.WeatherHandler)

	return mux
}

func startServer(mux *http.ServeMux, port string) error {
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Println("Starting server...")
	fmt.Printf("Running on http://localhost:%s\n", port)

	return server.ListenAndServe()
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error: .env file not found, using system variables")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS_URL string is required")
	}

	return &Config{
		RedisURL: redisURL,
	}, nil
}
