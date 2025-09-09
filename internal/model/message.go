package model

import "encoding/json"

// Message represents a flexible message structure that can handle
// both string and JSON body content
type Message struct {
	Headers    map[string]interface{} `json:"headers"`
	Exchange   string                 `json:"exchange"`
	RoutingKey string                 `json:"routingKey"`
	Body       json.RawMessage        `json:"body"`
}

// MarshalJSON custom marshaler to handle string or JSON body
func (m *Message) MarshalJSON() ([]byte, error) {
	// Create a temporary struct for marshaling
	msg := struct {
		Headers    map[string]interface{} `json:"headers"`
		Exchange   string                 `json:"exchange"`
		RoutingKey string                 `json:"routingKey"`
		Timestamp  int64                  `json:"timestamp"`
		Body       any                    `json:"body"`
	}{
		Headers:    m.Headers,
		Exchange:   m.Exchange,
		RoutingKey: m.RoutingKey,
	}

	// Try to unmarshal the body to detect if it's JSON or a string
	var bodyContent interface{}
	if err := json.Unmarshal(m.Body, &bodyContent); err == nil {
		msg.Body = bodyContent
	} else {
		msg.Body = string(m.Body)
	}

	return json.Marshal(msg)
}

// UnmarshalJSON custom unmarshaler to detect body type
func (m *Message) UnmarshalJSON(data []byte) error {
	var temp struct {
		Headers    map[string]interface{} `json:"headers"`
		Exchange   string                 `json:"exchange"`
		RoutingKey string                 `json:"routingKey"`
		Timestamp  int64                  `json:"timestamp"`
		Body       json.RawMessage        `json:"body"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Headers = temp.Headers
	m.Exchange = temp.Exchange
	m.RoutingKey = temp.RoutingKey

	var jsonCheck interface{}
	if err := json.Unmarshal(temp.Body, &jsonCheck); err == nil {
		m.Body = temp.Body
	} else {
		m.Body = temp.Body
	}

	return nil
}
