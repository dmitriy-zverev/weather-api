package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	apihandler "github.com/dmitriy-zverev/weather-api/internal/api_handler"
)

func BaseHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) WeatherHandler(w http.ResponseWriter, req *http.Request) {
	var params WeatherParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		answerWithError(w, "couldn't parse body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := apihandler.FetchCurrentWeather(
		context.Background(),
		cfg.RedisClient,
		params.City,
	)
	if err != nil {
		answerWithError(w, "couldn't fetch weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resData, err := json.Marshal(weather)
	if err != nil {
		answerWithError(w, "couldn't marshal weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resData)
}
