package rmq

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/streadway/amqp"
)

func (c *Consumer) filterMessage(msg *amqp.Delivery) bool {
	// Size filter
	if c.config.FilterConfig.MaxMessageSize > 0 && len(msg.Body) > c.config.FilterConfig.MaxMessageSize {
		return false
	}

	body := string(msg.Body)

	// Include patterns
	if len(c.config.FilterConfig.IncludePatterns) > 0 {
		if !containsAny(c.config.FilterConfig.IncludePatterns, body) {
			return false
		}
	}

	// Exclude patterns
	if containsAny(c.config.FilterConfig.ExcludePatterns, body) {
		return false
	}

	// JSON filter
	if c.config.FilterConfig.JSONFilter != nil {
		return matchJSONFilter(msg.Body, c.config.FilterConfig.JSONFilter)
	}

	return true
}

func (c *Consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func containsAny(patterns []string, body string) bool {
	for _, pattern := range patterns {
		if strings.Contains(body, pattern) {
			return true
		}
	}
	return false
}

func matchJSONFilter(body []byte, query *gojq.Query) bool {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false
	}

	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			fmt.Printf("Error applying JQ filter: %v\n", err)
			return false
		}

		// Check if result is truthy
		switch val := v.(type) {
		case bool:
			return val
		case nil:
			continue
		default:
			result, _ := json.MarshalIndent(v, "", "  ")
			return strings.TrimSpace(string(result)) != ""
		}
	}
	return false
}
