// Package weather provides clients, data models, formatters, and MCP server
// handlers for fetching and serving weather and geocoding information using
// the Open-Meteo API.
package weather

// Default coordinates and settings for Canberra, Australia.
const (
	// CanberraLatitude is the latitude for Canberra, Australia.
	CanberraLatitude = -35.2835
	// CanberraLongitude is the longitude for Canberra, Australia.
	CanberraLongitude = 149.1281
	// CanberraTimezone is the canonical timezone string for Canberra.
	CanberraTimezone = "Australia/Sydney"

	// DefaultForecastBaseURL is the base endpoint for the Open-Meteo Forecast API.
	DefaultForecastBaseURL = "https://api.open-meteo.com/v1/forecast"
	// DefaultGeocodingBaseURL is the base endpoint for the Open-Meteo Geocoding API.
	DefaultGeocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/search"
)

// CurrentWeather represents current weather measurements returned by Open-Meteo.
type CurrentWeather struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval,omitempty"`
	Temperature2m float64 `json:"temperature_2m"`
	WindSpeed10m  float64 `json:"wind_speed_10m"`
}

// DailyForecast represents aggregated daily weather metrics.
type DailyForecast struct {
	Time                        []string  `json:"time"`
	Temperature2mMax            []float64 `json:"temperature_2m_max"`
	Temperature2mMin            []float64 `json:"temperature_2m_min"`
	PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
}

// ForecastResponse represents the raw HTTP response payload from the Open-Meteo Forecast API.
type ForecastResponse struct {
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Timezone  string          `json:"timezone"`
	Elevation float64         `json:"elevation"`
	Current   *CurrentWeather `json:"current,omitempty"`
	Daily     *DailyForecast  `json:"daily,omitempty"`
}

// LocationResult represents a single geographic match from the Geocoding API.
type LocationResult struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
	Timezone  string  `json:"timezone"`
}

// GeocodingResponse represents the HTTP response payload from the Open-Meteo Geocoding API.
type GeocodingResponse struct {
	Results []LocationResult `json:"results"`
}
