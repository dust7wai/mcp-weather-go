package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLookupLocation_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "Canberra" {
			t.Errorf("expected query name 'Canberra', got '%s'", r.URL.Query().Get("name"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"id": 2172517,
					"name": "Canberra",
					"latitude": -35.28346,
					"longitude": 149.12807,
					"country": "Australia",
					"admin1": "Australian Capital Territory"
				}
			]
		}`))
	}))
	defer ts.Close()

	client := NewClient(WithGeocodingBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	loc, err := client.LookupLocation(context.Background(), "Canberra")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc == nil {
		t.Fatalf("expected location, got nil")
	}
	if loc.Name != "Canberra" || loc.Country != "Australia" || loc.Admin1 != "Australian Capital Territory" {
		t.Errorf("unexpected location fields: %+v", loc)
	}
}

func TestLookupLocation_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer ts.Close()

	client := NewClient(WithGeocodingBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	loc, err := client.LookupLocation(context.Background(), "NonExistentCity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc != nil {
		t.Fatalf("expected nil for not found, got %+v", loc)
	}
}

func TestLookupLocation_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer ts.Close()

	client := NewClient(WithGeocodingBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	_, err := client.LookupLocation(context.Background(), "Canberra")
	if err == nil {
		t.Fatalf("expected error on 500 status, got nil")
	}
}

func TestGetCurrentWeather_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"latitude": -35.2835,
			"longitude": 149.1281,
			"current": {
				"time": "2026-10-13T19:00",
				"temperature_2m": 13.9,
				"wind_speed_10m": 10.6
			}
		}`))
	}))
	defer ts.Close()

	client := NewClient(WithForecastBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	cw, err := client.GetCurrentWeather(context.Background(), CanberraLatitude, CanberraLongitude, CanberraTimezone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cw == nil {
		t.Fatalf("expected current weather, got nil")
	}
	if cw.Time != "2026-10-13T19:00" || cw.Temperature2m != 13.9 || cw.WindSpeed10m != 10.6 {
		t.Errorf("unexpected current weather values: %+v", cw)
	}
}

func TestGetDailyForecast_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"latitude": -35.2835,
			"longitude": 149.1281,
			"daily": {
				"time": ["2026-10-13", "2026-10-14"],
				"temperature_2m_max": [18.6, 16.8],
				"temperature_2m_min": [11.2, 6.0],
				"precipitation_probability_max": [84, 59]
			}
		}`))
	}))
	defer ts.Close()

	client := NewClient(WithForecastBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	df, err := client.GetDailyForecast(context.Background(), CanberraLatitude, CanberraLongitude, CanberraTimezone, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if df == nil {
		t.Fatalf("expected daily forecast, got nil")
	}
	if len(df.Time) != 2 || df.Temperature2mMax[0] != 18.6 || df.PrecipitationProbabilityMax[0] != 84 {
		t.Errorf("unexpected daily forecast: %+v", df)
	}
}

func TestGetCanberraHelpers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"current": {
				"time": "2026-10-13T19:00",
				"temperature_2m": 15.0,
				"wind_speed_10m": 5.0
			},
			"daily": {
				"time": ["2026-10-13"],
				"temperature_2m_max": [20.0],
				"temperature_2m_min": [10.0],
				"precipitation_probability_max": [10]
			}
		}`))
	}))
	defer ts.Close()

	client := NewClient(WithForecastBaseURL(ts.URL), WithHTTPClient(ts.Client()))

	cw, err := client.GetCanberraWeather(context.Background())
	if err != nil || cw == nil {
		t.Fatalf("GetCanberraWeather failed: %v", err)
	}

	df, err := client.GetCanberra7Days(context.Background())
	if err != nil || df == nil {
		t.Fatalf("GetCanberra7Days failed: %v", err)
	}
}

func TestGetCityWeather(t *testing.T) {
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("name") == "Tokyo" {
			_, _ = w.Write([]byte(`{
				"results": [{"id": 1, "name": "Tokyo", "latitude": 35.6895, "longitude": 139.6917, "country": "Japan"}]
			}`))
		} else {
			_, _ = w.Write([]byte(`{"results": []}`))
		}
	}))
	defer geoServer.Close()

	forecastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"current": {"time": "2026-10-13T19:00", "temperature_2m": 18.0, "wind_speed_10m": 8.0}
		}`))
	}))
	defer forecastServer.Close()

	client := NewClient(
		WithGeocodingBaseURL(geoServer.URL),
		WithForecastBaseURL(forecastServer.URL),
		WithHTTPClient(geoServer.Client()),
	)

	t.Run("found city", func(t *testing.T) {
		cw, loc, err := client.GetCityWeather(context.Background(), "Tokyo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cw == nil || loc == nil {
			t.Fatalf("expected cw and loc, got cw=%v loc=%v", cw, loc)
		}
		if loc.Name != "Tokyo" || cw.Temperature2m != 18.0 {
			t.Errorf("unexpected data: loc=%+v cw=%+v", loc, cw)
		}
	})

	t.Run("not found city", func(t *testing.T) {
		cw, loc, err := client.GetCityWeather(context.Background(), "Nowhere")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cw != nil || loc != nil {
			t.Fatalf("expected nil for not found city, got cw=%v loc=%v", cw, loc)
		}
	})
}
