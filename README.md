# GoDoctor - MCP Documentation Server

GoDoctor is a Model Context Protocol (MCP) server that provides tools for interacting with the Go environment. It exposes a documentation tool that allows AI models and users to retrieve Go package and symbol documentation directly from the terminal using the `go doc` command.

## Features

- **Model Context Protocol (MCP) Support:** Implements the official Go MCP SDK for seamless integration with MCP clients.
- **Documentation Retrieval:** A `read_docs` tool that invokes `go doc` to fetch documentation for any Go package or symbol.
- **Hello World Tool:** A simple `hello_world` tool to verify the server is running correctly.
- **CLI Client:** A dedicated test client for listing and calling tools from the command line.
- **Dockerized Environment:** Fully containerized development and execution environment.

---

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## For Users

### Running the Server

The server is designed to be run as an MCP server over `stdio` transport. You can launch it through the Docker container:

```bash
docker compose up -d --build
```

### Interacting with the Server (using the Test Client)

We provide a test client located at `./bin/client`. You can use it to list available tools and call them.

#### List Available Tools
```bash
docker compose exec app ./bin/client --tools-list
```

#### Call the Hello World Tool
```bash
docker compose exec app ./bin/client --tool-call hello_world
```

#### Retrieve Documentation for a Package
```bash
docker compose exec app ./bin/client --tool-call read_docs fmt
```

#### Retrieve Documentation for a Specific Symbol
```bash
docker compose exec app ./bin/client --tool-call read_docs fmt.Println
```

#### Retrieve Documentation for a Remote Package
```bash
docker compose exec app ./bin/client --tool-call read_docs github.com/modelcontextprotocol/go-sdk/mcp
```

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
    │       └── docs/      # Documentation retrieval logic
    ├── go.mod             # Go module definition
    └── GEMINI.md          # Development guidelines
```

### Building the Project

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
