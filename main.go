package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmitriy-zverev/weather-api/internal/handlers"
)

func main() {
	mux := setupRoutes()

	if err := startServer(mux, PORT); err != nil {
		log.Fatal("server error:", err)
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+BASE_ROUTE, handlers.BaseHandler)
	mux.HandleFunc("GET "+WEATHER_ROUTE, handlers.WeatherHandler)

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
