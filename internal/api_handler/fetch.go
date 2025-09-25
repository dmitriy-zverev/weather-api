package apihandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dmitriy-zverev/weather-api/internal/cache"
	"github.com/redis/go-redis/v9"
)

func FetchCurrentWeather(ctx context.Context, rdClient redis.Client, city string) (CurrentWeather, error) {
	weather, err := getCacheCurrentWeather(ctx, rdClient, city)
	if err == nil {
		return weather, nil
	}
	if err.Error() != NO_CACHE_DATA_ERROR && err != nil {
		log.Println("caching error:", err)
	}

	cityData, err := requestData(
		fmt.Sprintf("%s?name=%s&count=10&language=en&format=json&timezone=%s", GEO_URL, city, TIMEZONE),
	)
	if err != nil {
		return CurrentWeather{}, err
	}

	var geo geoAPI
	if err := json.Unmarshal(cityData, &geo); err != nil {
		return CurrentWeather{}, err
	}

	if err := cacheNewCity(ctx, rdClient, city, geo); err != nil {
		log.Println("caching city error: ", err)
	}

	weatherData, err := requestData(
		fmt.Sprintf(
			"%s?latitude=%.2f&longitude=%.2f&current=temperature_2m&timezone=%s",
			API_URL,
			geo.Results[0].Latitude,
			geo.Results[0].Longitude,
			TIMEZONE,
		),
	)
	if err != nil {
		return CurrentWeather{}, err
	}

	if err := json.Unmarshal(weatherData, &weather); err != nil {
		return CurrentWeather{}, err
	}

	if err := cacheNewCurrentTemp(ctx, rdClient, city, weather); err != nil {
		log.Println("caching current temperature error: ", err)
	}

	return weather, nil
}

func FetchWeatherNDays(ctx context.Context, rdClient redis.Client, city string, days int) (DailyWeather, error) {
	if days < 0 || days > 16 {
		return DailyWeather{}, errors.New("invalid number of days")
	}

	weatherDaily, err := getCacheNDaysWeather(ctx, rdClient, city, days)
	if err == nil {
		return weatherDaily, nil
	}
	if err.Error() != NO_CACHE_DATA_ERROR && err != nil {
		log.Println("caching error:", err)
	}

	cityData, err := requestData(
		fmt.Sprintf("%s?name=%s&count=10&language=en&format=json&timezone=%s", GEO_URL, city, TIMEZONE),
	)
	if err != nil {
		return DailyWeather{}, err
	}

	var geo geoAPI
	if err := json.Unmarshal(cityData, &geo); err != nil {
		return DailyWeather{}, err
	}

	if err := cacheNewCity(ctx, rdClient, city, geo); err != nil {
		log.Println("caching city error: ", err)
	}

	weatherDailyData, err := requestData(
		fmt.Sprintf(
			"%s?latitude=%.2f&longitude=%.2f&daily=temperature_2m_max,temperature_2m_min&timezone=%s&forecast_days=%d",
			API_URL,
			geo.Results[0].Latitude,
			geo.Results[0].Longitude,
			TIMEZONE,
			days,
		),
	)
	if err != nil {
		return DailyWeather{}, err
	}

	if err := json.Unmarshal(weatherDailyData, &weatherDaily); err != nil {
		return DailyWeather{}, err
	}

	if err := cacheNewNDaysTemp(ctx, rdClient, city, weatherDaily, days); err != nil {
		log.Println("caching current temperature error: ", err)
	}

	return weatherDaily, nil
}

func FetchWeatherHourly(ctx context.Context, rdClient redis.Client, city string) (HourlyWeather, error) {
	weatherHourly, err := getCacheHourlyWeather(ctx, rdClient, city)
	if err == nil {
		return weatherHourly, nil
	}
	if err.Error() != NO_CACHE_DATA_ERROR && err != nil {
		log.Println("caching error:", err)
	}

	cityData, err := requestData(
		fmt.Sprintf("%s?name=%s&count=10&language=en&format=json&timezone=%s", GEO_URL, city, TIMEZONE),
	)
	if err != nil {
		return HourlyWeather{}, err
	}

	var geo geoAPI
	if err := json.Unmarshal(cityData, &geo); err != nil {
		return HourlyWeather{}, err
	}

	if err := cacheNewCity(ctx, rdClient, city, geo); err != nil {
		log.Println("caching city error: ", err)
	}

	weatherHourlyData, err := requestData(
		fmt.Sprintf(
			"%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m&timezone=%s",
			API_URL,
			geo.Results[0].Latitude,
			geo.Results[0].Longitude,
			TIMEZONE,
		),
	)
	if err != nil {
		return HourlyWeather{}, err
	}

	if err := json.Unmarshal(weatherHourlyData, &weatherHourly); err != nil {
		return HourlyWeather{}, err
	}

	if err := cacheNewHourlyTemp(ctx, rdClient, city, weatherHourly); err != nil {
		log.Println("caching current temperature error: ", err)
	}

	return weatherHourly, nil
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

func getCacheNDaysWeather(ctx context.Context, client redis.Client, city string, days int) (DailyWeather, error) {
	cachedCityCoord, err := cache.Get(
		ctx,
		client,
		cache.GetCityKey(city),
	)
	if err != redis.Nil && err != nil {
		log.Println("redis error: ", err)
	}
	if err == nil && cachedCityCoord != "" {
		cachedData, err := cache.Get(
			ctx,
			client,
			cache.GetNDaysWeatherKey(cachedCityCoord, days),
		)
		if err != redis.Nil && err != nil {
			log.Println("redis error: ", err)
		}

		weatherData := []byte(cachedData)
		if cachedData == "" {
			weatherData, err = requestData(
				fmt.Sprintf(
					"%s?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min&timezone=%s&forecast_days=%d",
					API_URL,
					cachedCityCoord[:4],
					cachedCityCoord[5:],
					TIMEZONE,
					days,
				),
			)
			if err != nil {
				return DailyWeather{}, err
			}
		}

		var weatherDaily DailyWeather
		if err := json.Unmarshal(weatherData, &weatherDaily); err != nil {
			return DailyWeather{}, err
		}

		if err := cacheNewNDaysTemp(ctx, client, city, weatherDaily, days); err != nil {
			log.Println("caching current temperature error: ", err)
		}

		return weatherDaily, nil
	}

	return DailyWeather{}, errors.New(NO_CACHE_DATA_ERROR)
}

func getCacheHourlyWeather(ctx context.Context, client redis.Client, city string) (HourlyWeather, error) {
	cachedCityCoord, err := cache.Get(
		ctx,
		client,
		cache.GetCityKey(city),
	)
	if err != redis.Nil && err != nil {
		log.Println("redis error: ", err)
	}
	if err == nil && cachedCityCoord != "" {
		cachedData, err := cache.Get(
			ctx,
			client,
			cache.GetHourlyWeatherKey(cachedCityCoord),
		)
		if err != redis.Nil && err != nil {
			log.Println("redis error: ", err)
		}

		weatherData := []byte(cachedData)
		if cachedData == "" {
			weatherData, err = requestData(
				fmt.Sprintf(
					"%s?latitude=%s&longitude=%s&hourly=temperature_2m&timezone=%s",
					API_URL,
					cachedCityCoord[:4],
					cachedCityCoord[5:],
					TIMEZONE,
				),
			)
			if err != nil {
				return HourlyWeather{}, err
			}
		}

		var weatherHourly HourlyWeather
		if err := json.Unmarshal(weatherData, &weatherHourly); err != nil {
			return HourlyWeather{}, err
		}

		if err := cacheNewHourlyTemp(ctx, client, city, weatherHourly); err != nil {
			log.Println("caching current temperature error: ", err)
		}

		return weatherHourly, nil
	}

	return HourlyWeather{}, errors.New(NO_CACHE_DATA_ERROR)
}

func getCacheCurrentWeather(ctx context.Context, client redis.Client, city string) (CurrentWeather, error) {
	cachedCityCoord, err := cache.Get(
		ctx,
		client,
		cache.GetCityKey(city),
	)
	if err != redis.Nil && err != nil {
		log.Println("redis error: ", err)
	}
	if err == nil && cachedCityCoord != "" {
		cachedData, err := cache.Get(
			ctx,
			client,
			cache.GetCurrentWeatherKey(cachedCityCoord),
		)
		if err != redis.Nil && err != nil {
			log.Println("redis error: ", err)
		}

		weatherData := []byte(cachedData)
		if cachedData == "" {
			weatherData, err = requestData(
				fmt.Sprintf(
					"%s?latitude=%s&longitude=%s&current=temperature_2m&timezone=%s&format=json",
					API_URL,
					cachedCityCoord[:4],
					cachedCityCoord[5:],
					TIMEZONE,
				),
			)
			if err != nil {
				return CurrentWeather{}, err
			}
		}

		var weather CurrentWeather
		if err := json.Unmarshal(weatherData, &weather); err != nil {
			return CurrentWeather{}, err
		}

		if err := cacheNewCurrentTemp(ctx, client, city, weather); err != nil {
			log.Println("caching current temperature error: ", err)
		}

		return weather, nil
	}

	return CurrentWeather{}, errors.New(NO_CACHE_DATA_ERROR)
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

func cacheNewCurrentTemp(ctx context.Context, client redis.Client, city string, weather CurrentWeather) error {
	ttl, err := time.ParseDuration(TTL_CURRENT_TEMP)
	if err != nil {
		return err
	}

	cityCoord, err := cache.Get(ctx, client, cache.GetCityKey(city))
	if err != nil {
		return err
	}

	weatherData, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	if err := cache.Set(
		ctx,
		client,
		cache.GetCurrentWeatherKey(cityCoord),
		weatherData,
		ttl,
	); err != nil {
		return err
	}

	return nil
}

func cacheNewNDaysTemp(ctx context.Context, client redis.Client, city string, weatherDaily DailyWeather, days int) error {
	ttl, err := time.ParseDuration(TTL_DAILY_TEMP)
	if err != nil {
		return err
	}

	cityCoord, err := cache.Get(ctx, client, cache.GetCityKey(city))
	if err != nil {
		return err
	}

	weatherDailyData, err := json.Marshal(weatherDaily)
	if err != nil {
		return err
	}

	if err := cache.Set(
		ctx,
		client,
		cache.GetNDaysWeatherKey(cityCoord, days),
		weatherDailyData,
		ttl,
	); err != nil {
		return err
	}

	return nil
}

func cacheNewHourlyTemp(ctx context.Context, client redis.Client, city string, weatherHourly HourlyWeather) error {
	ttl, err := time.ParseDuration(TTL_DAILY_TEMP)
	if err != nil {
		return err
	}

	cityCoord, err := cache.Get(ctx, client, cache.GetCityKey(city))
	if err != nil {
		return err
	}

	weatherHourlyData, err := json.Marshal(weatherHourly)
	if err != nil {
		return err
	}

	if err := cache.Set(
		ctx,
		client,
		cache.GetHourlyWeatherKey(cityCoord),
		weatherHourlyData,
		ttl,
	); err != nil {
		return err
	}

	return nil
}
