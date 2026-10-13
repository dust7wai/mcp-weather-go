package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ClientOption defines a functional option for configuring a Client.
type ClientOption func(*Client)

// WithHTTPClient configures a custom *http.Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithForecastBaseURL overrides the default forecast API base endpoint.
func WithForecastBaseURL(rawURL string) ClientOption {
	return func(c *Client) {
		if rawURL != "" {
			c.forecastBaseURL = rawURL
		}
	}
}

// WithGeocodingBaseURL overrides the default geocoding API base endpoint.
func WithGeocodingBaseURL(rawURL string) ClientOption {
	return func(c *Client) {
		if rawURL != "" {
			c.geocodingBaseURL = rawURL
		}
	}
}

// WithLogger configures a structured logger for the client.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// Client interacts with the Open-Meteo Weather and Geocoding APIs.
type Client struct {
	forecastBaseURL  string
	geocodingBaseURL string
	httpClient       *http.Client
	logger           *slog.Logger
}

// NewClient creates and initializes a new Client instance.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		forecastBaseURL:  DefaultForecastBaseURL,
		geocodingBaseURL: DefaultGeocodingBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// LookupLocation looks up geographic coordinates for a given city or place name.
// Returns nil, nil if no location match is found.
func (c *Client) LookupLocation(ctx context.Context, name string) (*LocationResult, error) {
	reqURL := fmt.Sprintf("%s?name=%s&count=1&language=en",
		c.geocodingBaseURL,
		url.QueryEscape(name),
	)

	c.logger.DebugContext(ctx, "looking up location", "name", name, "url", reqURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("geocoding API returned status %d: %s", resp.StatusCode, string(body))
	}

	var geoResp GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return nil, fmt.Errorf("decode geocoding response failed: %w", err)
	}

	if len(geoResp.Results) == 0 {
		return nil, nil
	}

	return &geoResp.Results[0], nil
}

// GetCurrentWeather retrieves current weather measurements for given coordinates.
func (c *Client) GetCurrentWeather(ctx context.Context, lat, lon float64, timezone string) (*CurrentWeather, error) {
	if timezone == "" {
		timezone = "auto"
	}

	reqURL := fmt.Sprintf("%s?latitude=%v&longitude=%v&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m&timezone=%s",
		c.forecastBaseURL,
		lat,
		lon,
		url.QueryEscape(timezone),
	)

	c.logger.DebugContext(ctx, "fetching current weather", "lat", lat, "lon", lon, "timezone", timezone, "url", reqURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forecast API returned status %d: %s", resp.StatusCode, string(body))
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("decode forecast response failed: %w", err)
	}

	if forecast.Current == nil {
		return nil, fmt.Errorf("forecast API response missing current weather data")
	}

	return forecast.Current, nil
}

// GetDailyForecast retrieves the daily forecast for given coordinates.
func (c *Client) GetDailyForecast(ctx context.Context, lat, lon float64, timezone string, days int) (*DailyForecast, error) {
	if timezone == "" {
		timezone = "auto"
	}
	if days <= 0 {
		days = 7
	}

	reqURL := fmt.Sprintf("%s?latitude=%v&longitude=%v&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_max&timezone=%s&forecast_days=%d",
		c.forecastBaseURL,
		lat,
		lon,
		url.QueryEscape(timezone),
		days,
	)

	c.logger.DebugContext(ctx, "fetching daily forecast", "lat", lat, "lon", lon, "days", days, "url", reqURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forecast API returned status %d: %s", resp.StatusCode, string(body))
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("decode daily forecast response failed: %w", err)
	}

	if forecast.Daily == nil {
		return nil, fmt.Errorf("forecast API response missing daily forecast data")
	}

	return forecast.Daily, nil
}

// GetCanberraWeather retrieves current weather specifically for Canberra, Australia.
func (c *Client) GetCanberraWeather(ctx context.Context) (*CurrentWeather, error) {
	return c.GetCurrentWeather(ctx, CanberraLatitude, CanberraLongitude, CanberraTimezone)
}

// GetCanberra7Days retrieves the 7-day weather forecast for Canberra, Australia.
func (c *Client) GetCanberra7Days(ctx context.Context) (*DailyForecast, error) {
	return c.GetDailyForecast(ctx, CanberraLatitude, CanberraLongitude, CanberraTimezone, 7)
}

// GetCityWeather looks up the city coordinates and fetches current weather for that location.
// Returns (nil, nil, nil) if the city was not found in the geocoding service.
func (c *Client) GetCityWeather(ctx context.Context, city string) (*CurrentWeather, *LocationResult, error) {
	loc, err := c.LookupLocation(ctx, city)
	if err != nil {
		return nil, nil, fmt.Errorf("location lookup failed: %w", err)
	}
	if loc == nil {
		return nil, nil, nil
	}

	cw, err := c.GetCurrentWeather(ctx, loc.Latitude, loc.Longitude, "auto")
	if err != nil {
		return nil, loc, fmt.Errorf("fetch weather failed: %w", err)
	}

	return cw, loc, nil
}
