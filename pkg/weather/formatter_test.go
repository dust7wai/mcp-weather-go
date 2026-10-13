package weather

import (
	"strings"
	"testing"
)

func TestFormatCanberraWeather(t *testing.T) {
	cw := &CurrentWeather{
		Time:          "2026-10-13T19:00",
		Temperature2m: 13.9,
		WindSpeed10m:  10.6,
	}

	formatted := FormatCanberraWeather(cw)

	expected := "☀️ Canberra Current Weather ☀️\nTime:    2026-10-13T19:00\nTemperature: 13.9°C\nWind speed: 10.6km/h"
	if formatted != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, formatted)
	}

	if FormatCanberraWeather(nil) != "" {
		t.Errorf("expected empty string for nil input")
	}
}

func TestFormatCityWeather(t *testing.T) {
	cw := &CurrentWeather{
		Time:          "2026-10-13T12:00",
		Temperature2m: 21.5,
		WindSpeed10m:  15.2,
	}

	formatted := FormatCityWeather(cw)

	expected := "🌍 Weather Report 🌍\nTime:    2026-10-13T12:00\nTemperature: 21.5°C\nWind speed: 15.2km/h"
	if formatted != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, formatted)
	}

	if FormatCityWeather(nil) != "" {
		t.Errorf("expected empty string for nil input")
	}
}

func TestFormatLocation(t *testing.T) {
	t.Run("with admin1", func(t *testing.T) {
		loc := &LocationResult{
			Name:      "Canberra",
			Admin1:    "Australian Capital Territory",
			Country:   "Australia",
			Latitude:  -35.28346,
			Longitude: 149.12807,
		}

		formatted := FormatLocation(loc)
		expected := "📍 Found Location: Canberra, Australian Capital Territory, Australia\nLatitude: -35.28346\nLongitude: 149.12807"
		if formatted != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s", expected, formatted)
		}
	})

	t.Run("without admin1", func(t *testing.T) {
		loc := &LocationResult{
			Name:      "Singapore",
			Admin1:    "",
			Country:   "Singapore",
			Latitude:  1.3521,
			Longitude: 103.8198,
		}

		formatted := FormatLocation(loc)
		expected := "📍 Found Location: Singapore, Singapore\nLatitude: 1.3521\nLongitude: 103.8198"
		if formatted != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s", expected, formatted)
		}
	})

	t.Run("empty country", func(t *testing.T) {
		loc := &LocationResult{
			Name:      "Unknown Place",
			Admin1:    "",
			Country:   "",
			Latitude:  0.0,
			Longitude: 0.0,
		}

		formatted := FormatLocation(loc)
		expected := "📍 Found Location: Unknown Place, Unknown Country\nLatitude: 0\nLongitude: 0"
		if formatted != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s", expected, formatted)
		}
	})

	t.Run("nil location", func(t *testing.T) {
		if FormatLocation(nil) != "" {
			t.Errorf("expected empty string for nil input")
		}
	})
}

func TestFormatDailyForecast(t *testing.T) {
	df := &DailyForecast{
		Time:                        []string{"2026-10-13", "2026-10-14"},
		Temperature2mMax:            []float64{18.6, 16.8},
		Temperature2mMin:            []float64{11.2, 6.0},
		PrecipitationProbabilityMax: []float64{84, 59},
	}

	formatted := FormatDailyForecast(df)

	if !strings.Contains(formatted, "☀️ Canberra 7-day Weather Forecast ☀️") {
		t.Errorf("missing header in formatted output")
	}
	if !strings.Contains(formatted, "Date        Day         Min     Max     Rain  ") {
		t.Errorf("missing table column headers")
	}
	if !strings.Contains(formatted, "2026-10-13  Tuesday      11.2°C  18.6°C   84%") {
		t.Errorf("missing or mismatched row 1 in:\n%s", formatted)
	}
	if !strings.Contains(formatted, "2026-10-14  Wednesday     6.0°C  16.8°C   59%") {
		t.Errorf("missing or mismatched row 2 in:\n%s", formatted)
	}

	if FormatDailyForecast(nil) != "" {
		t.Errorf("expected empty string for nil input")
	}
}
