// Package main provides the primary entrypoint for the mcp-weather-go application.
// By default, it launches the Model Context Protocol (MCP) server over standard I/O.
// It also provides subcommands for displaying Canberra's current weather or 7-day forecast.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dust7wai/mcp-weather-go/pkg/weather"
)

func printHelp() {
	fmt.Println("mcp-weather-go - Model Context Protocol Weather Server & CLI Tools in Go")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run . [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  server (default)  Launch the MCP server over stdio")
	fmt.Println("  now               Fetch and display current Canberra weather")
	fmt.Println("  7days             Fetch and display Canberra 7-day weather forecast")
	fmt.Println("  help              Show this help message")
}

func main() {
	cmd := "server"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	client := weather.NewClient()

	switch cmd {
	case "server":
		srv := weather.NewServer(client)
		if err := weather.ServeStdio(srv); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server encountered error: %v\n", err)
			os.Exit(1)
		}

	case "now", "current":
		cw, err := client.GetCanberraWeather(context.Background())
		if err != nil {
			fmt.Printf("Error fetching weather data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(weather.FormatCanberraWeather(cw))

	case "7days", "forecast":
		df, err := client.GetCanberra7Days(context.Background())
		if err != nil {
			fmt.Printf("Error fetching weather data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(weather.FormatDailyForecast(df))

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}
