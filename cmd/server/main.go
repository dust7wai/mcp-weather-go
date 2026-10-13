// Package main provides the entry point for running the Model Context Protocol (MCP) weather server over stdio.
package main

import (
	"fmt"
	"os"

	"github.com/dust7wai/mcp-weather-go/pkg/weather"
)

func main() {
	client := weather.NewClient()
	srv := weather.NewServer(client)

	if err := weather.ServeStdio(srv); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server encountered error: %v\n", err)
		os.Exit(1)
	}
}
