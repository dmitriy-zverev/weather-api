package apihandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dmitriy-zverev/weather-api/internal/cache"
	"github.com/redis/go-redis/v9"
)

func FetchCurrentWeather(ctx context.Context, rdClient redis.Client, city string) (Weather, error) {
	cachedCityCoord, err := cache.Get(
		ctx,
		rdClient,
		cache.GetCityKey(city),
	)
	if err != nil {
		log.Println("redis error: ", err)
	}
	if err == nil && cachedCityCoord != "" {
		cachedData, err := cache.Get(
			ctx,
			rdClient,
			cache.GetCurrentWeatherKey(cachedCityCoord),
		)
		if err != nil {
			log.Println("redis error: ", err)
		}

		weatherData := []byte(cachedData)
		if cachedData == "" {
			weatherData, err = requestData(
				fmt.Sprintf(
					"%s?latitude=%s&longitude=%s&current=temperature_2m",
					API_URL,
					cachedCityCoord[:3],
					cachedCityCoord[3:],
				),
			)
			if err != nil {
				return Weather{}, err
			}
		}

		var weather Weather
		if err := json.Unmarshal(weatherData, &weather); err != nil {
			return Weather{}, err
		}

		return weather, nil
	}

	data, err := requestData(
		fmt.Sprintf("%s?name=%s&count=10&language=en&format=json", GEO_URL, city),
	)
	if err != nil {
		return Weather{}, err
	}

	var geo geoAPI
	if err := json.Unmarshal(data, &geo); err != nil {
		return Weather{}, err
	}

	if err := cacheNewCity(ctx, rdClient, city, geo); err != nil {
		log.Println("caching city error: ", err)
	}

	weatherData, err := requestData(
		fmt.Sprintf(
			"%s?latitude=%.2f&longitude=%.2f&current=temperature_2m",
			API_URL,
			geo.Results[0].Latitude,
			geo.Results[0].Longitude,
		),
	)
	if err != nil {
		return Weather{}, err
	}

	var weather Weather
	if err := json.Unmarshal(weatherData, &weather); err != nil {
		return Weather{}, err
	}

	if err := cacheNewTemp(ctx, rdClient, geo, weather); err != nil {
		log.Println("caching current temperature error: ", err)
	}

	return weather, nil
}

func requestData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func cacheNewCity(ctx context.Context, client redis.Client, cityName string, geo geoAPI) error {
	ttl, err := time.ParseDuration(TTL_CITY)
	if err != nil {
		return err
	}

	coords := fmt.Sprintf("%.2f", geo.Results[0].Latitude) + "," + fmt.Sprintf("%.2f", geo.Results[0].Longitude)

	if err := cache.Set(
		ctx,
		client,
		cache.GetCityKey(cityName),
		[]byte(coords),
		ttl,
	); err != nil {
		return err
	}

	return nil
}

func cacheNewTemp(ctx context.Context, client redis.Client, geo geoAPI, weather Weather) error {
	ttl, err := time.ParseDuration(TTL_CURRENT_TEMP)
	if err != nil {
		return err
	}

	coords := fmt.Sprintf("%.2f", geo.Results[0].Latitude) + "," + fmt.Sprintf("%.2f", geo.Results[0].Longitude)

	weatherData, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	if err := cache.Set(
		ctx,
		client,
		cache.GetCurrentWeatherKey(coords),
		weatherData,
		ttl,
	); err != nil {
		return err
	}

	return nil
}
