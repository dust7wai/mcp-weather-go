package weather

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// FormatCanberraWeather formats current weather metrics for Canberra.
func FormatCanberraWeather(cw *CurrentWeather) string {
	if cw == nil {
		return ""
	}
	return fmt.Sprintf("☀️ Canberra Current Weather ☀️\nTime:    %s\nTemperature: %v°C\nWind speed: %vkm/h",
		cw.Time,
		cw.Temperature2m,
		cw.WindSpeed10m,
	)
}

// FormatCityWeather formats current weather metrics for an arbitrary city.
func FormatCityWeather(cw *CurrentWeather) string {
	if cw == nil {
		return ""
	}
	return fmt.Sprintf("🌍 Weather Report 🌍\nTime:    %s\nTemperature: %v°C\nWind speed: %vkm/h",
		cw.Time,
		cw.Temperature2m,
		cw.WindSpeed10m,
	)
}

// FormatLocation formats geographic metadata for a discovered location.
func FormatLocation(loc *LocationResult) string {
	if loc == nil {
		return ""
	}

	country := loc.Country
	if country == "" {
		country = "Unknown Country"
	}

	var regionInfo string
	if loc.Admin1 != "" {
		regionInfo = fmt.Sprintf("%s, %s, %s", loc.Name, loc.Admin1, country)
	} else {
		regionInfo = fmt.Sprintf("%s, %s", loc.Name, country)
	}

	return fmt.Sprintf("📍 Found Location: %s\nLatitude: %v\nLongitude: %v",
		regionInfo,
		loc.Latitude,
		loc.Longitude,
	)
}

// FormatDailyForecast formats a 7-day daily weather forecast table.
func FormatDailyForecast(df *DailyForecast) string {
	if df == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("☀️ Canberra 7-day Weather Forecast ☀️\n\n")
	sb.WriteString(fmt.Sprintf("%-12s%-12s%-8s%-8s%-6s\n", "Date", "Day", "Min", "Max", "Rain"))

	count := len(df.Time)
	for i := 0; i < count; i++ {
		dateStr := df.Time[i]
		dayName := "Unknown"
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			dayName = parsedDate.Weekday().String()
		}

		minTemp := 0.0
		if i < len(df.Temperature2mMin) {
			minTemp = df.Temperature2mMin[i]
		}

		maxTemp := 0.0
		if i < len(df.Temperature2mMax) {
			maxTemp = df.Temperature2mMax[i]
		}

		rainProb := 0.0
		if i < len(df.PrecipitationProbabilityMax) {
			rainProb = df.PrecipitationProbabilityMax[i]
		}

		rainInt := int(math.Round(rainProb))

		sb.WriteString(fmt.Sprintf("%-12s%-12s%5.1f°C %5.1f°C %4d%%\n",
			dateStr,
			dayName,
			minTemp,
			maxTemp,
			rainInt,
		))
	}

	return strings.TrimRight(sb.String(), "\n")
}
