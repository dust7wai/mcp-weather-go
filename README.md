# mcp-weather-go

A Model Context Protocol (MCP) server and command-line weather tools implemented in Go.

Powered by the [Open-Meteo API](https://open-meteo.com/).

---

## Features

| Feature | Go Implementation | Description |
|---|---|---|
| **Current Weather CLI** | `cmd/weathernow` / `go run . now` | Fetches and displays Canberra's current temperature and wind speed. |
| **7-Day Forecast CLI** | `cmd/weather7days` / `go run . 7days` | Displays a 7-day Canberra forecast table (Max, Min, Rain probability). |
| **MCP Server (Stdio)** | `cmd/server` / `go run .` | MCP server exposing tools over standard input/output. |

### MCP Tools Provided

1. `lookup_location`: Looks up geographic coordinates (latitude and longitude) for a given city or place name using Open-Meteo Geocoding.
2. `get_weather`: Retrieves current weather, temperature, and wind speed for any global city name.
3. `get_canberra_weather`: Retrieves current weather, temperature, and wind speed specifically for Canberra, Australia.

---

## Project Structure

```text
mcp-weather-go/
├── cmd/
│   ├── server/             # MCP server entrypoint
│   │   └── main.go
│   ├── weathernow/         # Current Canberra weather CLI
│   │   └── main.go
│   └── weather7days/       # Canberra 7-day forecast CLI
│       └── main.go
├── pkg/
│   └── weather/            # Reusable weather package
│       ├── client.go       # Open-Meteo HTTP client
│       ├── client_test.go  # Client unit tests (httptest)
│       ├── formatter.go    # Output formatters (matching Python output)
│       ├── formatter_test.go
│       ├── models.go       # Data types and request/response structures
│       ├── server.go       # MCP server definition & tool handlers
│       └── server_test.go  # Server & tool execution unit tests
├── Makefile                # Build, test, and run automation
├── main.go                 # Unified CLI and default MCP server runner
├── go.mod
├── go.sum
└── README.md
```

---

## Requirements

- **Go 1.22+** (tested with Go 1.27)

---

## Quick Start

### 1. Run via `go run .` (Unified CLI)

```bash
# Run the MCP server (default)
go run .

# Fetch current Canberra weather
go run . now

# Fetch Canberra 7-day weather forecast
go run . 7days

# View CLI help
go run . help
```

### 2. Run Individual Command Binaries

```bash
# Current Canberra weather
go run ./cmd/weathernow

# 7-day weather forecast
go run ./cmd/weather7days

# MCP Server
go run ./cmd/server
```

### 3. Using `make`

```bash
# Build all binaries into bin/
make build

# Run unit tests
make test

# Format and vet code
make fmt
make vet

# Execute tools
make run-now
make run-7days
make run-server
```

---

## MCP Client Configuration

To register this server with an MCP client (such as Claude Desktop, Cursor, or Antigravity), add the following entry to your MCP configuration file:

### Using `go run` directly:

```json
{
  "mcpServers": {
    "weather": {
      "command": "go",
      "args": [
        "run",
        "/path/to/mcp-weather-go/cmd/server"
      ]
    }
  }
}
```

### Using the compiled binary:

```bash
make build
```

```json
{
  "mcpServers": {
    "weather": {
      "command": "/path/to/mcp-weather-go/bin/weather-server"
    }
  }
}
```

---

## Testing

Run the full test suite:

```bash
go test -v ./...
```

Unit tests use `net/http/httptest` to mock Open-Meteo API responses, ensuring fast, deterministic, offline execution.
