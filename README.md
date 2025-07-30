# MCP Go Development Container

This is a development container setup for building Model Context Protocol (MCP) servers using Go.

## Features

- **Go 1.22** development environment
- **Pre-configured VS Code** with Go extensions
- **MCP server implementation** with basic tools
- **WebSocket support** for MCP communication
- **Health monitoring** endpoint
- **Hot reload** with Air (for development)
- **Static analysis** tools (golangci-lint, staticcheck)

## Getting Started

### Prerequisites

- Docker
- VS Code with Dev Containers extension

### Development Setup

1. **Open in Dev Container**
   - Open VS Code in this folder
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
   - Select "Dev Containers: Reopen in Container"

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the Server**
   ```bash
   go run main.go
   ```

4. **Test the Server**
   - Health check: `curl http://localhost:8080/health`
   - WebSocket endpoint: `ws://localhost:8080/mcp`

## MCP Server Features

### Available Tools

1. **Echo Tool**
   - Echoes back any message
   - Useful for testing MCP communication

2. **Timestamp Tool**
   - Returns the current timestamp
   - Demonstrates parameter-less tools

### MCP Protocol Support

- Protocol version: `2024-11-05`
- JSON-RPC 2.0 communication
- WebSocket transport
- Tool execution capabilities

## Development Commands

### Run with Hot Reload
```bash
air
```

### Lint Code
```bash
golangci-lint run
```

### Run Static Analysis
```bash
staticcheck ./...
```

### Format Code
```bash
goimports -w .
```

### Build for Production
```bash
go build -o mcp-server main.go
```

## Testing the MCP Server

### Using WebSocket Client

You can test the MCP server using a WebSocket client. Here's an example initialization message:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}
```

### List Available Tools

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### Call Echo Tool

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello, MCP!"
    }
  }
}
```

## Project Structure

```
.
├── .devcontainer/
│   ├── devcontainer.json    # Dev container configuration
│   └── Dockerfile          # Custom container setup
├── static/                 # Static files for testing
├── main.go                # MCP server implementation
├── go.mod                 # Go module definition
└── README.md              # This file
```

## Adding New Tools

To add a new tool to the MCP server:

1. Add the tool definition in the `handleToolsList` method
2. Add the tool execution logic in the `handleToolCall` method
3. Test the new tool using the WebSocket interface

## Environment Variables

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (default: info)

## Contributing

1. Make your changes
2. Run tests: `go test ./...`
3. Lint code: `golangci-lint run`
4. Format code: `goimports -w .`

## Resources

- [Model Context Protocol Specification](https://spec.modelcontextprotocol.io/)
- [Go Documentation](https://golang.org/doc/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)