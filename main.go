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
	"golang.org/x/time/rate"
)

type Config struct {
	RedisURL string
	Port     string
	Platform string
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

	fmt.Println("Starting Redis...")
	fmt.Printf("Redis running on http://%s\n\n", config.RedisURL)

	// 7 requests per minute (based on API free usage)
	limiter := rate.NewLimiter(0.1, 10)

	apiConfig := &handlers.ApiConfig{
		RedisClient: *rdb,
		Limiter:     limiter,
	}

	mux := setupRoutes(apiConfig)

	if err := startServer(mux, PORT); err != nil {
		log.Fatal("server error:", err)
	}
}

func setupRoutes(cfg *handlers.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+API_VERSION+BASE_ROUTE, handlers.BaseHandler)
	mux.HandleFunc("GET "+API_VERSION+WEATHER_ROUTE, cfg.WeatherHandler)

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

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		return nil, errors.New("PLATFORM string is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT string is required")
	}

	rdPort := os.Getenv("REDIS_PORT")
	if rdPort == "" {
		return nil, errors.New("REDIS_PORT string is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS_URL string is required")
	}
	if platform == "prod" {
		redisURL = "redis"
	}

	return &Config{
		RedisURL: redisURL + ":" + rdPort,
		Port:     port,
		Platform: platform,
	}, nil
}
