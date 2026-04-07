# GoDoctor - MCP Documentation & Code Review Server

GoDoctor is a Model Context Protocol (MCP) server that provides tools for interacting with the Go environment. It exposes a documentation tool, a code review tool powered by Gemini on Vertex AI, and more.

## Features

| Tool | Description |
|------|-------------|
| `hello_world` | Returns a hello world message to verify the server is running |
| `read_docs` | Invokes `go doc` to fetch documentation for any Go package or symbol |
| `code_review` | Analyzes Go code using Gemini 2.5 Pro on Vertex AI and returns improvements in Markdown |
| `sbom_generate` | Parses `go.mod` and generates a Software Bill of Materials in Markdown |
| `go_test` | Runs `go test` on a specified package and returns the results |
| `go_vet` | Runs `go vet` for static analysis on a specified package |

Additional features:

- **Model Context Protocol (MCP) Support:** Implements the official Go MCP SDK for seamless integration with MCP clients.
- **Streamable HTTP Transport:** Supports both stdio and HTTP transport for Cloud Run deployment.
- **CLI Client:** A dedicated test client for listing and calling tools from the command line.
- **Dockerized Environment:** Fully containerized development and execution environment.

---

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+ (or Docker + Docker Compose)
- A Google Cloud project with Vertex AI API enabled (for the `code_review` tool)

---

## For Users

### Running the Server

The server is designed to be run as an MCP server over `stdio` transport.

#### With Docker
```bash
docker compose up -d --build
```

#### Without Docker
```bash
CGO_ENABLED=0 go build -o bin/server ./cmd/server/
CGO_ENABLED=0 go build -o bin/client ./cmd/client/
```

### Interacting with the Server (using the Test Client)

We provide a test client located at `./bin/client`. You can use it to list available tools and call them.

#### List Available Tools
```bash
# With Docker
docker compose exec app ./bin/client --tools-list

# Without Docker
./bin/client --tools-list
```

#### Call the Hello World Tool
```bash
docker compose exec app ./bin/client --tool-call hello_world
# or
./bin/client --tool-call hello_world
```

#### Retrieve Documentation for a Package
```bash
docker compose exec app ./bin/client --tool-call read_docs fmt
# or
./bin/client --tool-call read_docs fmt
```

#### Retrieve Documentation for a Specific Symbol
```bash
./bin/client --tool-call read_docs fmt.Println
```

#### Retrieve Documentation for a Remote Package
```bash
./bin/client --tool-call read_docs github.com/modelcontextprotocol/go-sdk/mcp
```

#### Review Go Code with AI
```bash
./bin/client --tool-call code_review internal/tools/code/review.go
```

#### Review Go Code with a Specific Focus
```bash
./bin/client --tool-call code_review --hint "focus on security" internal/tools/code/review.go
```

#### Generate SBOM from go.mod
```bash
./bin/client --tool-call sbom_generate go.mod
```

#### Run Tests on a Package
```bash
./bin/client --tool-call go_test ./internal/tools/...
```

#### Run Static Analysis on a Package
```bash
./bin/client --tool-call go_vet ./cmd/server/...
```

> **Note:** The `code_review` tool requires the server to be started with `--project` and `--location` flags for Vertex AI authentication. See the server configuration below.

---

## For Developers

### Project Structure

```text
/
├── Dockerfile             # Docker configuration
├── docker-compose.yml    # Docker Compose setup
└── godoctor/              # Main Go application
    ├── bin/               # Compiled binaries
    ├── cmd/               # Binary entry points
    │   ├── client/        # Test client implementation
    │   └── server/        # MCP server implementation
    ├── internal/          # Internal packages
    │   └── tools/         # MCP tool implementations
    │       ├── code/      # AI code review logic
    │       ├── docs/      # Documentation retrieval logic
    │       ├── gotest/    # Go test runner
    │       ├── govet/     # Go vet static analysis
    │       └── sbom/      # SBOM generation
    ├── go.mod             # Go module definition
    └── GEMINI.md          # Development guidelines
```

### Building the Project

#### With Docker

All builds must be performed inside the Docker container to ensure consistent environment and dependencies.

1. **Rebuild the container (if Dockerfile changed):**
   ```bash
   docker compose up -d --build
   ```

2. **Compile the binaries:**
   ```bash
   docker compose exec app go build -o bin/server cmd/server/main.go
   docker compose exec app go build -o bin/client cmd/client/main.go
   ```

#### Without Docker

Requires Go 1.25+ installed locally.

```bash
CGO_ENABLED=0 go build -o bin/server ./cmd/server/
CGO_ENABLED=0 go build -o bin/client ./cmd/client/
```

### Server Flags

| Flag | Description |
|------|-------------|
| `--project` | Google Cloud Project ID (required for `code_review`). Defaults to `GOOGLE_CLOUD_PROJECT` env var |
| `--location` | Google Cloud Location (required for `code_review`). Defaults to `GOOGLE_CLOUD_LOCATION` env var |
| `-http` | Use Streamable HTTP transport instead of stdio |
| `-listen` | Address to listen on with `-http` (default: `:8080`) |

### Client Flags

| Flag | Description |
|------|-------------|
| `--tools-list` | List available tools |
| `--tool-call` | Call a specific tool |
| `--hint` | Optional hint for `code_review` (e.g. `"focus on security"`) |
| `--server-path` | Path to the server binary (default: `./bin/server`) |
| `-addr` | Server HTTP endpoint (e.g. `http://localhost:8080/mcp`) |

### Adding a New Tool

To add a new tool to the MCP server:

1. Create a new package under `internal/tools/`.
2. Implement the tool's logic and define its arguments as a struct.
3. Register the tool in `cmd/server/main.go` using the `mcp.AddTool` function.

Example registration:
```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "new_tool",
    Description: "Description of what it does",
}, func(ctx context.Context, req *mcp.CallToolRequest, args MyToolArgs) (*mcp.CallToolResult, any, error) {
    // Implementation logic here
})
```

### Development Guidelines

Please refer to [GEMINI.md](./GEMINI.md) for detailed coding standards, including:
- Formatting with `gofmt`.
- Naming conventions (camelCase/PascalCase).
- Error handling patterns.
- Proper documentation for exported identifiers.
- Docker-based execution commands.

---

## License

This project follows the coding standards and principles outlined in the [Google Go Style Guide](https://google.github.io/styleguide/go/).
