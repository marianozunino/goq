package model

import (
	"encoding/json"
	"testing"
)

func TestMessage_MarshalJSON_RealImplementation(t *testing.T) {
	// Test the actual MarshalJSON method
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
			"priority":     5,
		},
		Exchange:   "test_exchange",
		RoutingKey: "test.key",
		Body:       json.RawMessage(`{"test": "data"}`),
	}

	// Marshal to JSON using the actual method
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

	// Compare JSON bodies by unmarshaling them
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

func TestMessage_MarshalJSON_EmptyMessage(t *testing.T) {
	// Test with minimal message
	msg := Message{}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal empty message: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty message: %v", err)
	}

	// Verify empty message
	if result.Exchange != "" {
		t.Error("Expected empty exchange")
	}

	if result.RoutingKey != "" {
		t.Error("Expected empty routing key")
	}

	// Body should be null JSON or empty
	if len(result.Body) > 0 && string(result.Body) != "null" {
		t.Errorf("Expected empty body or null, got length %d, body: %v", len(result.Body), result.Body)
	}

	// Headers should be nil or empty
	if result.Headers != nil && len(result.Headers) > 0 {
		t.Error("Expected nil or empty headers")
	}
}

func TestMessage_MarshalJSON_OnlyHeaders(t *testing.T) {
	// Test with only headers
	msg := Message{
		Headers: map[string]interface{}{
			"test": "value",
		},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with only headers: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with only headers: %v", err)
	}

	if result.Headers["test"] != "value" {
		t.Errorf("Expected header 'test' to be 'value', got %v", result.Headers["test"])
	}
}

func TestMessage_MarshalJSON_OnlyBody(t *testing.T) {
	// Test with only body
	msg := Message{
		Body: json.RawMessage(`{"data": "test"}`),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with only body: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with only body: %v", err)
	}

	// Compare JSON bodies by unmarshaling them
	var expectedBody, resultBody map[string]interface{}
	json.Unmarshal([]byte(`{"data": "test"}`), &expectedBody)
	json.Unmarshal(result.Body, &resultBody)

	if len(expectedBody) != len(resultBody) {
		t.Errorf("Expected body length %d, got %d", len(expectedBody), len(resultBody))
	}
}

func TestMessage_MarshalJSON_OnlyExchange(t *testing.T) {
	// Test with only exchange
	msg := Message{
		Exchange: "test_exchange",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with only exchange: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with only exchange: %v", err)
	}

	if result.Exchange != "test_exchange" {
		t.Errorf("Expected exchange to be 'test_exchange', got %s", result.Exchange)
	}
}

func TestMessage_MarshalJSON_OnlyRoutingKey(t *testing.T) {
	// Test with only routing key
	msg := Message{
		RoutingKey: "test.key",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with only routing key: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with only routing key: %v", err)
	}

	if result.RoutingKey != "test.key" {
		t.Errorf("Expected routing key to be 'test.key', got %s", result.RoutingKey)
	}
}

func TestMessage_MarshalJSON_AllFields(t *testing.T) {
	// Test with all fields populated
	msg := Message{
		Headers: map[string]interface{}{
			"content-type": "application/json",
			"priority":     5,
			"timestamp":    "2023-01-01T00:00:00Z",
		},
		Exchange:   "events",
		RoutingKey: "user.created",
		Body:       json.RawMessage(`{"user": {"id": 123, "name": "test"}}`),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal complete message: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal complete message: %v", err)
	}

	// Verify all fields
	if result.Exchange != "events" {
		t.Errorf("Expected exchange 'events', got %s", result.Exchange)
	}

	if result.RoutingKey != "user.created" {
		t.Errorf("Expected routing key 'user.created', got %s", result.RoutingKey)
	}

	// Compare JSON bodies by unmarshaling them
	var expectedBody, resultBody map[string]interface{}
	json.Unmarshal([]byte(`{"user": {"id": 123, "name": "test"}}`), &expectedBody)
	json.Unmarshal(result.Body, &resultBody)

	if len(expectedBody) != len(resultBody) {
		t.Errorf("Expected body length %d, got %d", len(expectedBody), len(resultBody))
	}

	if len(result.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(result.Headers))
	}

	if result.Headers["content-type"] != "application/json" {
		t.Errorf("Expected content-type 'application/json', got %v", result.Headers["content-type"])
	}

	if result.Headers["priority"] != float64(5) {
		t.Errorf("Expected priority 5, got %v", result.Headers["priority"])
	}

	if result.Headers["timestamp"] != "2023-01-01T00:00:00Z" {
		t.Errorf("Expected timestamp '2023-01-01T00:00:00Z', got %v", result.Headers["timestamp"])
	}
}

func TestMessage_MarshalJSON_InvalidJSONBody(t *testing.T) {
	// Test with invalid JSON in body - this should fail during marshaling
	msg := Message{
		Body: json.RawMessage(`{invalid json`),
	}

	// This should fail because json.RawMessage expects valid JSON
	_, err := json.Marshal(msg)
	if err == nil {
		t.Error("Expected error when marshaling invalid JSON body")
	}
}

func TestMessage_MarshalJSON_NilBody(t *testing.T) {
	// Test with nil body
	msg := Message{
		Body: nil,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message with nil body: %v", err)
	}

	var result Message
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal message with nil body: %v", err)
	}

	// Body should be null JSON or empty
	if len(result.Body) > 0 && string(result.Body) != "null" {
		t.Errorf("Expected empty body or null, got length %d, body: %v", len(result.Body), result.Body)
	}
}
