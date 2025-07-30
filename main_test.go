package main

import (
	"encoding/json"
	"testing"
)

func TestMCPMessage(t *testing.T) {
	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
		Params:  map[string]interface{}{"key": "value"},
	}

	// Test JSON marshaling
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal MCPMessage: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled MCPMessage
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPMessage: %v", err)
	}

	if unmarshaled.JSONRPC != msg.JSONRPC {
		t.Errorf("Expected JSONRPC %s, got %s", msg.JSONRPC, unmarshaled.JSONRPC)
	}

	if unmarshaled.Method != msg.Method {
		t.Errorf("Expected Method %s, got %s", msg.Method, unmarshaled.Method)
	}
}

func TestMCPServer_handleInitialize(t *testing.T) {
	server := NewMCPServer()

	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
		},
	}

	response := server.handleInitialize(msg)

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC 2.0, got %s", response.JSONRPC)
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %v", response.ID)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	protocolVersion, ok := result["protocolVersion"].(string)
	if !ok || protocolVersion != "2024-11-05" {
		t.Errorf("Expected protocolVersion 2024-11-05, got %v", protocolVersion)
	}
}

func TestMCPServer_handleToolsList(t *testing.T) {
	server := NewMCPServer()

	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	response := server.handleToolsList(msg)

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	tools, ok := result["tools"].([]Tool)
	if !ok {
		t.Fatal("Expected tools to be a slice of Tool")
	}

	if len(tools) == 0 {
		t.Error("Expected at least one tool")
	}

	// Check for echo tool
	foundEcho := false
	for _, tool := range tools {
		if tool.Name == "echo" {
			foundEcho = true
			break
		}
	}

	if !foundEcho {
		t.Error("Expected to find echo tool")
	}
}

func TestMCPServer_handleToolCall_Echo(t *testing.T) {
	server := NewMCPServer()

	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "echo",
			"arguments": map[string]interface{}{
				"message": "test message",
			},
		},
	}

	response := server.handleToolCall(msg)

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	content, ok := result["content"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected content to be a slice of maps")
	}

	if len(content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(content))
	}

	text, ok := content[0]["text"].(string)
	if !ok {
		t.Fatal("Expected text to be a string")
	}

	expected := "Echo: test message"
	if text != expected {
		t.Errorf("Expected text %s, got %s", expected, text)
	}
}

func TestMCPServer_handleToolCall_Timestamp(t *testing.T) {
	server := NewMCPServer()

	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      "timestamp",
			"arguments": map[string]interface{}{},
		},
	}

	response := server.handleToolCall(msg)

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	content, ok := result["content"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected content to be a slice of maps")
	}

	if len(content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(content))
	}

	text, ok := content[0]["text"].(string)
	if !ok {
		t.Fatal("Expected text to be a string")
	}

	if text == "" {
		t.Error("Expected non-empty timestamp text")
	}
}
