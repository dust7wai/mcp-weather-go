// Package main provides a CLI command to retrieve and display current Canberra weather.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dust7wai/mcp-weather-go/pkg/weather"
)

// RunWeatherNow retrieves and prints current Canberra weather from the Open-Meteo API.
func RunWeatherNow(ctx context.Context) error {
	client := weather.NewClient()
	cw, err := client.GetCanberraWeather(ctx)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		return err
	}

	fmt.Println(weather.FormatCanberraWeather(cw))
	return nil
}

// WeatherNow is an alias for RunWeatherNow for backward compatibility.
var WeatherNow = RunWeatherNow

func main() {
	if err := RunWeatherNow(context.Background()); err != nil {
		os.Exit(1)
	}
}
