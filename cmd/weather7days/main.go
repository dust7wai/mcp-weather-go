// Package main provides a CLI command to retrieve and display Canberra's 7-day forecast.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dust7wai/mcp-weather-go/pkg/weather"
)

// RunWeather7Days retrieves and prints the 7-day weather forecast for Canberra.
func RunWeather7Days(ctx context.Context) error {
	client := weather.NewClient()
	df, err := client.GetCanberra7Days(ctx)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		return err
	}

	fmt.Println(weather.FormatDailyForecast(df))
	return nil
}

func main() {
	if err := RunWeather7Days(context.Background()); err != nil {
		os.Exit(1)
	}
}
