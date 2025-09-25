package cache

import (
	"fmt"
	"strings"
)

const (
	NAMESPACE   = "wx"
	API_VERSION = "v1"
)

func GetCityKey(cityName string) string {
	return fmt.Sprintf(
		"%s:%s:city:by_name:%s",
		NAMESPACE,
		API_VERSION,
		strings.ToLower(cityName),
	)
}

func GetCurrentWeatherKey(cityCoord string) string {
	return fmt.Sprintf(
		"%s:%s:weather:latest:%s",
		NAMESPACE,
		API_VERSION,
		strings.ToLower(cityCoord),
	)
}
