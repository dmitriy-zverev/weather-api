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

func GetNDaysWeatherKey(cityCoord string, days int) string {
	return fmt.Sprintf(
		"%s:%s:weather:%d_days:%s",
		NAMESPACE,
		API_VERSION,
		days,
		strings.ToLower(cityCoord),
	)
}

func GetHourlyWeatherKey(cityCoord string) string {
	return fmt.Sprintf(
		"%s:%s:weather:hourly:%s",
		NAMESPACE,
		API_VERSION,
		strings.ToLower(cityCoord),
	)
}
