package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// MCPServer represents the Model Context Protocol server
type MCPServer struct {
	logger   *logrus.Logger
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

// MCPMessage represents a basic MCP message structure
type MCPMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error in MCP format
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &MCPServer{
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// handleWebSocket handles WebSocket connections
func (s *MCPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
		return
	}
	defer conn.Close()

	s.clients[conn] = true
	defer delete(s.clients, conn)

	s.logger.Info("Client connected")

	for {
		var msg MCPMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			s.logger.WithError(err).Error("Failed to read message")
			break
		}

		s.logger.WithField("message", msg).Info("Received message")

		response := s.handleMessage(msg)
		if response != nil {
			err = conn.WriteJSON(response)
			if err != nil {
				s.logger.WithError(err).Error("Failed to send response")
				break
			}
		}
	}

	s.logger.Info("Client disconnected")
}

// handleMessage processes MCP messages
func (s *MCPServer) handleMessage(msg MCPMessage) *MCPMessage {
	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			},
		}
	}
}

// handleInitialize handles the initialize request
func (s *MCPServer) handleInitialize(msg MCPMessage) *MCPMessage {
	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "mcp-go-server",
				"version": "1.0.0",
			},
		},
	}
}

// handleToolsList returns available tools
func (s *MCPServer) handleToolsList(msg MCPMessage) *MCPMessage {
	tools := []Tool{
		{
			Name:        "echo",
			Description: "Echo back the input message",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "The message to echo back",
					},
				},
				"required": []string{"message"},
			},
		},
		{
			Name:        "timestamp",
			Description: "Get the current timestamp",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		// TODO: Add new tool definitions here
		// Example:
		// {
		//     Name:        "your_tool_name",
		//     Description: "Description of your tool",
		//     InputSchema: map[string]interface{}{
		//         "type": "object",
		//         "properties": map[string]interface{}{
		//             "param_name": map[string]interface{}{
		//                 "type":        "string",
		//                 "description": "Parameter description",
		//             },
		//         },
		//         "required": []string{"param_name"},
		//     },
		// },
	}

	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

// handleToolCall executes a tool
func (s *MCPServer) handleToolCall(msg MCPMessage) *MCPMessage {
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	name, ok := params["name"].(string)
	if !ok {
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Tool name is required",
			},
		}
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	var result interface{}
	var err error

	switch name {
	case "echo":
		message, _ := arguments["message"].(string)
		result = map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Echo: %s", message),
				},
			},
		}
	case "timestamp":
		result = map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Current timestamp: %s", time.Now().Format(time.RFC3339)),
				},
			},
		}
	default:
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", name),
			},
		}
	}

	if err != nil {
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: fmt.Sprintf("Tool execution failed: %v", err),
			},
		}
	}

	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  result,
	}
}

// health endpoint for monitoring
func (s *MCPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"clients":   len(s.clients),
	})
}

func main() {
	server := NewMCPServer()

	router := mux.NewRouter()
	router.HandleFunc("/mcp", server.handleWebSocket)
	router.HandleFunc("/health", server.handleHealth).Methods("GET")

	// Serve static files for testing
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		server.logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			server.logger.WithError(err).Error("Server shutdown failed")
		}
	}()

	server.logger.WithField("address", srv.Addr).Info("Starting MCP server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	server.logger.Info("Server stopped")
}
