package apihandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchCurrentWeather(city string) (Weather, error) {
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
