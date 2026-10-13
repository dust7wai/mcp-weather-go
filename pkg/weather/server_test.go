package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewServer_ToolsRegistration(t *testing.T) {
	client := NewClient()
	s := NewServer(client)

	tools := s.ListTools()
	expectedTools := []string{"lookup_location", "get_weather", "get_canberra_weather"}

	for _, name := range expectedTools {
		tool, ok := tools[name]
		if !ok || tool == nil {
			t.Errorf("expected tool %s to be registered", name)
		}
	}
}

func TestServerToolHandlers(t *testing.T) {
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		name := r.URL.Query().Get("name")
		if name == "Canberra" {
			_, _ = w.Write([]byte(`{
				"results": [{
					"id": 1,
					"name": "Canberra",
					"latitude": -35.2835,
					"longitude": 149.1281,
					"country": "Australia",
					"admin1": "ACT"
				}]
			}`))
		} else {
			_, _ = w.Write([]byte(`{"results": []}`))
		}
	}))
	defer geoServer.Close()

	forecastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"current": {
				"time": "2026-10-13T19:00",
				"temperature_2m": 14.5,
				"wind_speed_10m": 12.0
			}
		}`))
	}))
	defer forecastServer.Close()

	client := NewClient(
		WithGeocodingBaseURL(geoServer.URL),
		WithForecastBaseURL(forecastServer.URL),
		WithHTTPClient(geoServer.Client()),
	)

	s := NewServer(client)

	t.Run("lookup_location found", func(t *testing.T) {
		st := s.GetTool("lookup_location")
		req := mcp.CallToolRequest{}
		req.Params.Arguments = map[string]any{"name": "Canberra"}

		res, err := st.Handler(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected handler error: %v", err)
		}
		text := res.Content[0].(mcp.TextContent).Text
		if !strings.Contains(text, "📍 Found Location: Canberra, ACT, Australia") {
			t.Errorf("unexpected output: %s", text)
		}
	})

	t.Run("lookup_location not found", func(t *testing.T) {
		st := s.GetTool("lookup_location")
		req := mcp.CallToolRequest{}
		req.Params.Arguments = map[string]any{"name": "MissingPlace"}

		res, err := st.Handler(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected handler error: %v", err)
		}
		text := res.Content[0].(mcp.TextContent).Text
		if !strings.Contains(text, "❌ Could not find any location coordinates matching 'MissingPlace'.") {
			t.Errorf("unexpected output: %s", text)
		}
	})

	t.Run("get_weather found", func(t *testing.T) {
		st := s.GetTool("get_weather")
		req := mcp.CallToolRequest{}
		req.Params.Arguments = map[string]any{"city": "Canberra"}

		res, err := st.Handler(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected handler error: %v", err)
		}
		text := res.Content[0].(mcp.TextContent).Text
		if !strings.Contains(text, "🌍 Weather Report 🌍") || !strings.Contains(text, "14.5°C") {
			t.Errorf("unexpected output: %s", text)
		}
	})

	t.Run("get_canberra_weather", func(t *testing.T) {
		st := s.GetTool("get_canberra_weather")
		req := mcp.CallToolRequest{}

		res, err := st.Handler(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected handler error: %v", err)
		}
		text := res.Content[0].(mcp.TextContent).Text
		if !strings.Contains(text, "☀️ Canberra Current Weather ☀️") || !strings.Contains(text, "14.5°C") {
			t.Errorf("unexpected output: %s", text)
		}
	})
}
