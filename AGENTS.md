# Go Development Guidelines
All code contributed to this project must adhere to the following principles.

## 1. Formatting & Linting
All Go code **must** be formatted and linted using official Go tooling before being submitted. 
- Use **`go fmt`** or **`goimports`** strictly to handle indentation, layout, and automatic import sorting.
- Run **`golangci-lint`** or **`go vet`** to catch static analysis bugs, performance anti-patterns, and shadowed variables before declaring a feature complete.

## 2. Naming Conventions
- **Packages:** Use short, all-lowercase, single-word names. Do **NOT** use underscores, dashes, or `mixedCaps` in package names (e.g., use `weatherserver`, not `weather_server`).
- **Exported Identifiers (Public):** Use `PascalCase` starting with a capital letter (e.g., `McpServer`). This determines visibility outside the package.
- **Unexported Identifiers (Private):** Use `camelCase` starting with a lowercase letter (e.g., `lookupLocation`).
- **Interfaces:** Usually name interfaces with an "-er" suffix if they only contain one or two functions (e.g., `Reader`, `Formatter`).
- **Acronyms:** Keep acronyms uniformly capitalized. Use `ServeHTTP` or `apiURL`, not `ServeHttp` or `apiUrl`.

## 3. Error Handling
- **Explicit Returns:** Go does not use exceptions. Functions that can fail must return an `error` as their final value. Always inspect this error immediately:
  ```go
  result, err := lookupLocation(ctx, input)
  if err != nil {
      return fmt.Errorf("lookup location failed: %w", err)
  }
  ```
- **Error Wrapping:** Use `fmt.Errorf("... %w", err)` to wrap contextual details around lower-level failures without hiding the underlying root cause.
- **Stream Protection:** Do **NOT** use standard output (`fmt.Println()`) for general logging or debug prints. This pollutes the JSON-RPC stream used for `stdio` transport and will crash the connection. **Always use standard structured logging (`slog`) directed to `os.Stderr`.**

## 4. Simplicity and Clarity
- Write clean, flat, idiomatic code. Guard clauses with early returns are highly preferred over deeply nested `if-else` blocks.
- Minimize interface abstractions. "Accept interfaces, return structs." Do not over-engineer architectures with unnecessary abstract layers.
- Avoid naked or named returns in long or complex functions, as they obscure readability.

## 5. Documentation
- All public packages, exported structs, interfaces, methods, and functions **must** have descriptive comments adhering to Godoc conventions.
- Comments for exported symbols must start with the name of the symbol itself (e.g., `// McpServer represents the primary runtime engine...`).
- Tool and handler documentation is structurally critical, as LLMs read these text signatures to understand functionality. Ensure argument boundaries and requirements are explicit.

## 6. Project Structure & Dependency Management
- Follow the standard Go project layout:
  - `/cmd` - Application entry points (e.g., `/cmd/server/main.go`).
  - `/pkg` - Shared code safe for external applications to import.
  - `/internal` - Private application code that the Go compiler explicitly blocks external projects from importing.
- **Dependency Management:** Use Go Modules. The project root must contain a valid `go.mod` and a locked `go.sum` file.
- **Executables:** Use standard module commands at the root: `go run`, `go test ./...`, and `go build`.

## 7. Executable Core Commands
- **Install/Tidy Dependencies:** `go mod tidy`
- **Run Unit Tests:** `go test -v ./...`
- **Run Code Formatter:** `go fmt ./...`
- **Run Complete Linter Suite:** `golangci-lint run`
- **Execute Local Binary:** `go run cmd/server/main.go`
