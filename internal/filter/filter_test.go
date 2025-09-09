package filter

import (
	"testing"

	"github.com/marianozunino/goq/internal/config"
)

// testMessage implements MessageDelivery interface for testing
type testMessage struct {
	body []byte
}

func (t *testMessage) GetBody() []byte {
	return t.body
}

func TestMessageFilter_NoFilters(t *testing.T) {
	cfg := &config.Config{}
	filter := NewMessageFilter(cfg)

	if errs := filter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Unexpected compilation errors: %v", errs)
	}

	msg := &testMessage{body: []byte(`{"test": "data"}`)}
	if !filter.Filter(msg) {
		t.Error("Expected message to pass filter with no filters configured")
	}
}

func TestMessageFilter_IncludePattern(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{"admin"},
			ExcludePatterns: []string{},
			JSONFilter:      "",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	if errs := filter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Unexpected compilation errors: %v", errs)
	}

	// Should pass - contains "admin"
	msg1 := &testMessage{body: []byte(`{"user": "admin", "action": "login"}`)}
	if !filter.Filter(msg1) {
		t.Error("Expected message with 'admin' to pass include filter")
	}

	// Should fail - doesn't contain "admin"
	msg2 := &testMessage{body: []byte(`{"user": "user", "action": "login"}`)}
	if filter.Filter(msg2) {
		t.Error("Expected message without 'admin' to fail include filter")
	}
}

func TestMessageFilter_ExcludePattern(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{"error"},
			JSONFilter:      "",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	if errs := filter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Unexpected compilation errors: %v", errs)
	}

	// Should fail - contains "error"
	msg1 := &testMessage{body: []byte(`{"status": "error", "message": "failed"}`)}
	if filter.Filter(msg1) {
		t.Error("Expected message with 'error' to fail exclude filter")
	}

	// Should pass - doesn't contain "error"
	msg2 := &testMessage{body: []byte(`{"status": "success", "message": "ok"}`)}
	if !filter.Filter(msg2) {
		t.Error("Expected message without 'error' to pass exclude filter")
	}
}

func TestMessageFilter_JSONFilter(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			JSONFilter:      `.user.role == "admin"`,
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	if errs := filter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Unexpected compilation errors: %v", errs)
	}

	// Should pass - user role is admin
	msg1 := &testMessage{body: []byte(`{"user": {"role": "admin", "name": "john"}}`)}
	if !filter.Filter(msg1) {
		t.Error("Expected message with admin role to pass JSON filter")
	}

	// Should fail - user role is not admin
	msg2 := &testMessage{body: []byte(`{"user": {"role": "user", "name": "jane"}}`)}
	if filter.Filter(msg2) {
		t.Error("Expected message with non-admin role to fail JSON filter")
	}
}

func TestMessageFilter_MaxMessageSize(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			JSONFilter:      "",
			MaxMessageSize:  10,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	if errs := filter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Unexpected compilation errors: %v", errs)
	}

	// Should pass - message is small enough
	msg1 := &testMessage{body: []byte(`{"a": 1}`)}
	if !filter.Filter(msg1) {
		t.Error("Expected small message to pass size filter")
	}

	// Should fail - message is too large
	msg2 := &testMessage{body: []byte(`{"very": "long message that exceeds the limit"}`)}
	if filter.Filter(msg2) {
		t.Error("Expected large message to fail size filter")
	}
}

func TestMessageFilter_InvalidRegexPattern(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{"[invalid"},
			ExcludePatterns: []string{},
			JSONFilter:      "",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	errs := filter.GetCompilationErrors()
	if len(errs) == 0 {
		t.Error("Expected compilation errors for invalid regex pattern")
	}

	// Should not panic when filtering with invalid pattern
	msg := &testMessage{body: []byte(`{"test": "data"}`)}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Filter panicked with invalid regex: %v", r)
		}
	}()

	filter.Filter(msg)
}

func TestMessageFilter_InvalidJSONFilter(t *testing.T) {
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			JSONFilter:      "invalid jq syntax",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := NewMessageFilter(cfg)
	errs := filter.GetCompilationErrors()
	if len(errs) == 0 {
		t.Error("Expected compilation errors for invalid JSON filter")
	}
}
