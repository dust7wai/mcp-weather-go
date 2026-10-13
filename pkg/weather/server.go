package weather

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ServerName is the default identification name for the MCP server.
const ServerName = "Global Weather & Location Server"

// ServerVersion is the semantic version string for the MCP server.
const ServerVersion = "1.0.0"

// NewServer creates a new MCP server populated with weather and location tools.
func NewServer(client *Client, opts ...server.ServerOption) *server.MCPServer {
	if client == nil {
		client = NewClient()
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	defaultOpts := []server.ServerOption{
		server.WithToolCapabilities(true),
		server.WithLogger(logger),
	}
	allOpts := append(defaultOpts, opts...)

	s := server.NewMCPServer(ServerName, ServerVersion, allOpts...)

	// Tool 1: lookup_location
	lookupTool := mcp.NewTool("lookup_location",
		mcp.WithDescription("Looks up geographic coordinates (latitude and longitude) for a given city or place name."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the city or place to look up (e.g., 'Canberra', 'Paris', 'Tokyo')."),
		),
	)
	s.AddTool(lookupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid argument 'name': %v", err)), nil
		}

		loc, err := client.LookupLocation(ctx, name)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error during location lookup for '%s': %v", name, err)), nil
		}
		if loc == nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Could not find any location coordinates matching '%s'.", name)), nil
		}

		return mcp.NewToolResultText(FormatLocation(loc)), nil
	})

	// Tool 2: get_weather
	weatherTool := mcp.NewTool("get_weather",
		mcp.WithDescription("Retrieves the current weather, temperature, and wind speed for any global location using coordinates."),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("The name of the city or place to look up (e.g., 'Canberra', 'Paris', 'Tokyo')."),
		),
	)
	s.AddTool(weatherTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid argument 'city': %v", err)), nil
		}

		loc, err := client.LookupLocation(ctx, city)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error fetching weather data for coordinates (<nil>, <nil>) of city %s: %v", city, err)), nil
		}
		if loc == nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Could not find any location coordinates matching '%s'.", city)), nil
		}

		cw, err := client.GetCurrentWeather(ctx, loc.Latitude, loc.Longitude, "auto")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error fetching weather data for coordinates (%v, %v) of city %s: %v", loc.Latitude, loc.Longitude, city, err)), nil
		}

		return mcp.NewToolResultText(FormatCityWeather(cw)), nil
	})

	// Tool 3: get_canberra_weather
	canberraTool := mcp.NewTool("get_canberra_weather",
		mcp.WithDescription("Retrieves the current weather, temperature, and wind speed for Canberra, Australia."),
	)
	s.AddTool(canberraTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cw, err := client.GetCanberraWeather(ctx)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error fetching weather data: %v", err)), nil
		}

		return mcp.NewToolResultText(FormatCanberraWeather(cw)), nil
	})

	return s
}

// ServeStdio starts the MCP server over standard input and standard output.
func ServeStdio(s *server.MCPServer) error {
	errorLogger := log.New(os.Stderr, "[mcp-weather-server] ", log.LstdFlags)
	if err := server.ServeStdio(s, server.WithErrorLogger(errorLogger)); err != nil {
		return fmt.Errorf("stdio server failed: %w", err)
	}
	return nil
}
