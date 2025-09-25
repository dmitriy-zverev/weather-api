package apihandler

type geoAPI struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type Weather struct {
	Current struct {
		Temp float64 `json:"temperature_2m"`
	} `json:"current"`
}
