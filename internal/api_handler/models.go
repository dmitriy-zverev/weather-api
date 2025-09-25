package apihandler

type geoAPI struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type CurrentWeather struct {
	Current struct {
		Temp float64 `json:"temperature_2m"`
	} `json:"current"`
}

type DailyWeather struct {
	Daily struct {
		Time    []string  `json:"time"`
		MaxTemp []float64 `json:"temperature_2m_max"`
		MinTemp []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type HourlyWeather struct {
	Hourly struct {
		Time []string  `json:"time"`
		Temp []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}
