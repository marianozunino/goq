package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMessage_MarshalJSON(t *testing.T) {
	// Create a test message
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
			"priority":     5,
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       json.RawMessage(`{"test": "data"}`),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify fields
	if result.Exchange != msg.Exchange {
		t.Errorf("Expected exchange %s, got %s", msg.Exchange, result.Exchange)
	}

	if result.RoutingKey != msg.RoutingKey {
		t.Errorf("Expected routing key %s, got %s", msg.RoutingKey, result.RoutingKey)
	}

	// Verify body (JSON formatting may differ, so we'll parse and compare)
	var originalBody, resultBody map[string]interface{}
	json.Unmarshal(msg.Body, &originalBody)
	json.Unmarshal(result.Body, &resultBody)

	if len(originalBody) != len(resultBody) {
		t.Errorf("Expected body length %d, got %d", len(originalBody), len(resultBody))
	}

	// Verify headers
	if len(result.Headers) != len(msg.Headers) {
		t.Errorf("Expected %d headers, got %d", len(msg.Headers), len(result.Headers))
	}

	if result.Headers["content-type"] != "application/json" {
		t.Errorf("Expected content-type 'application/json', got %v", result.Headers["content-type"])
	}

	if result.Headers["priority"] != float64(5) {
		t.Errorf("Expected priority 5, got %v", result.Headers["priority"])
	}
}

func TestMessage_MarshalJSON_EmptyHeaders(t *testing.T) {
	// Create a message with empty headers
	msg := Message{
		Headers:    nil,
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       []byte(`{"test": "data"}`),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with empty headers: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with empty headers: %v", err)
	}

	// Verify headers are nil
	if result.Headers != nil {
		t.Error("Expected headers to be nil")
	}
}

func TestMessage_MarshalJSON_EmptyBody(t *testing.T) {
	// Create a message with empty body
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       json.RawMessage(`{}`),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with empty body: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with empty body: %v", err)
	}

	// Verify body is empty object
	if string(result.Body) != "{}" {
		t.Errorf("Expected empty object body, got %s", string(result.Body))
	}
}

func TestMessage_MarshalJSON_ComplexHeaders(t *testing.T) {
	// Create a message with complex headers
	msg := Message{
		Headers: map[string]interface{}{
			"string-header": "test-value",
			"int-header":    42,
			"float-header":  3.14,
			"bool-header":   true,
			"array-header":  []interface{}{"item1", "item2"},
			"map-header": map[string]interface{}{
				"nested-key": "nested-value",
			},
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       []byte(`{"test": "data"}`),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with complex headers: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with complex headers: %v", err)
	}

	// Verify complex headers
	if result.Headers["string-header"] != "test-value" {
		t.Errorf("Expected string-header 'test-value', got %v", result.Headers["string-header"])
	}

	if result.Headers["int-header"] != float64(42) {
		t.Errorf("Expected int-header 42, got %v", result.Headers["int-header"])
	}

	if result.Headers["float-header"] != 3.14 {
		t.Errorf("Expected float-header 3.14, got %v", result.Headers["float-header"])
	}

	if result.Headers["bool-header"] != true {
		t.Errorf("Expected bool-header true, got %v", result.Headers["bool-header"])
	}

	// Verify array header
	arrayHeader, ok := result.Headers["array-header"].([]interface{})
	if !ok {
		t.Error("Expected array-header to be an array")
	} else if len(arrayHeader) != 2 {
		t.Errorf("Expected array-header to have 2 items, got %d", len(arrayHeader))
	}

	// Verify map header
	mapHeader, ok := result.Headers["map-header"].(map[string]interface{})
	if !ok {
		t.Error("Expected map-header to be a map")
	} else if mapHeader["nested-key"] != "nested-value" {
		t.Errorf("Expected nested-key 'nested-value', got %v", mapHeader["nested-key"])
	}
}

func TestMessage_MarshalJSON_BinaryBody(t *testing.T) {
	// Create a message with binary-like data encoded as base64 in JSON
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/octet-stream",
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       json.RawMessage(`{"data": "AAECAwP//w=="}`),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with binary body: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with binary body: %v", err)
	}

	// Verify binary body (JSON formatting may differ)
	var originalBody, resultBody map[string]interface{}
	json.Unmarshal(msg.Body, &originalBody)
	json.Unmarshal(result.Body, &resultBody)

	if len(originalBody) != len(resultBody) {
		t.Errorf("Expected body length %d, got %d", len(originalBody), len(resultBody))
	}
}

func TestMessage_MarshalJSON_LargeBody(t *testing.T) {
	// Create a message with large JSON body
	largeBody := json.RawMessage(`{"large": "` + strings.Repeat("x", 1000) + `"}`)
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       largeBody,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with large body: %v", err)
	}

	// Unmarshal back to verify
	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with large body: %v", err)
	}

	// Verify large body
	if len(result.Body) < 1000 {
		t.Errorf("Expected large body, got %d bytes", len(result.Body))
	}
}

func TestMessage_String(t *testing.T) {
	// Create a test message
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       []byte(`{"test": "data"}`),
	}

	// Test that we can create the message without panicking
	if msg.Exchange != "test_exchange" {
		t.Error("Expected exchange to be set correctly")
	}

	if msg.RoutingKey != "test.key" {
		t.Error("Expected routing key to be set correctly")
	}
}
