.PHONY: all build test fmt vet run-now run-7days run-server clean tidy

BINARY_DIR=bin

all: build

build:
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/weather-server ./cmd/server
	go build -o $(BINARY_DIR)/weather-now ./cmd/weathernow
	go build -o $(BINARY_DIR)/weather-7days ./cmd/weather7days

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

run-now:
	go run ./cmd/weathernow

run-7days:
	go run ./cmd/weather7days

run-server:
	go run ./cmd/server

clean:
	rm -rf $(BINARY_DIR)
