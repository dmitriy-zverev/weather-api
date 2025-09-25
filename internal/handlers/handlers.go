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

	switch params.ForecastType {
	case CURRENT:
		weather, err := apihandler.FetchCurrentWeather(
			context.Background(),
			cfg.RedisClient,
			params.City,
		)
		if err != nil {
			answerWithError(w, "couldn't fetch weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writeResponseBody(w, weather); err != nil {
			answerWithError(w, "couldn't marshal weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case N_DAYS:
		weatherDaily, err := apihandler.FetchWeatherNDays(
			context.Background(),
			cfg.RedisClient,
			params.City,
			params.Days,
		)
		if err != nil {
			answerWithError(w, "couldn't fetch weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writeResponseBody(w, weatherDaily); err != nil {
			answerWithError(w, "couldn't marshal weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case TODAY_HOURLY:
		weatherHourly, err := apihandler.FetchWeatherHourly(
			context.Background(),
			cfg.RedisClient,
			params.City,
		)
		if err != nil {
			answerWithError(w, "couldn't fetch weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writeResponseBody(w, weatherHourly); err != nil {
			answerWithError(w, "couldn't marshal weather data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		answerWithError(w, "you did not specify 'forecast_type'", http.StatusBadRequest)
		return
	}

}

func writeResponseBody[T any](w http.ResponseWriter, val T) error {
	resData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resData)

	return nil
}
