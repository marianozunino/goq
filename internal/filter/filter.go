package filter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/itchyny/gojq"
	"github.com/marianozunino/goq/internal/config"
)

type MessageFilter struct {
	maxMessageSize    int
	includePatterns   []*regexp.Regexp
	excludePatterns   []*regexp.Regexp
	jsonFilter        *gojq.Query
	regexFilter       *regexp.Regexp
	compilationErrors []error
	mu                sync.RWMutex
}

func NewMessageFilter(cfg *config.Config) *MessageFilter {
	filter := &MessageFilter{
		maxMessageSize: cfg.FilterConfig.MaxMessageSize,
	}

	// Compile regex patterns with error tracking
	filter.compilePatterns(cfg)

	return filter
}

func (f *MessageFilter) compilePatterns(cfg *config.Config) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.compilationErrors = []error{}

	// Compile include patterns
	f.includePatterns = make([]*regexp.Regexp, 0, len(cfg.FilterConfig.IncludePatterns))
	for _, pattern := range cfg.FilterConfig.IncludePatterns {
		if regex, err := regexp.Compile(pattern); err != nil {
			f.compilationErrors = append(f.compilationErrors, fmt.Errorf("invalid include pattern: %v", err))
		} else {
			f.includePatterns = append(f.includePatterns, regex)
		}
	}

	// Compile exclude patterns
	f.excludePatterns = make([]*regexp.Regexp, 0, len(cfg.FilterConfig.ExcludePatterns))
	for _, pattern := range cfg.FilterConfig.ExcludePatterns {
		if regex, err := regexp.Compile(pattern); err != nil {
			f.compilationErrors = append(f.compilationErrors, fmt.Errorf("invalid exclude pattern: %v", err))
		} else {
			f.excludePatterns = append(f.excludePatterns, regex)
		}
	}

	// Compile JSON filter
	if cfg.FilterConfig.JSONFilter != "" {
		query, err := gojq.Parse(cfg.FilterConfig.JSONFilter)
		if err != nil {
			f.compilationErrors = append(f.compilationErrors, fmt.Errorf("invalid JSON filter: %v", err))
		} else {
			f.jsonFilter = query
		}
	}

	// Compile regex filter
	if cfg.FilterConfig.RegexFilter != "" {
		if regex, err := regexp.Compile(cfg.FilterConfig.RegexFilter); err != nil {
			f.compilationErrors = append(f.compilationErrors, fmt.Errorf("invalid regex filter: %v", err))
		} else {
			f.regexFilter = regex
		}
	}
}

// MessageDelivery represents a message that can be filtered
type MessageDelivery interface {
	GetBody() []byte
}

func (f *MessageFilter) Filter(msg MessageDelivery) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Size filter
	if f.maxMessageSize > 0 && len(msg.GetBody()) > f.maxMessageSize {
		return false
	}

	body := string(msg.GetBody())

	// Regex filter
	if f.regexFilter != nil && !f.regexFilter.MatchString(body) {
		return false
	}

	// Include patterns
	if len(f.includePatterns) > 0 {
		if !f.matchAnyRegex(f.includePatterns, body) {
			return false
		}
	}

	// Exclude patterns
	if f.matchAnyRegex(f.excludePatterns, body) {
		return false
	}

	// JSON filter
	if f.jsonFilter != nil {
		return f.matchJSONFilter(msg.GetBody())
	}

	return true
}

func (f *MessageFilter) matchAnyRegex(patterns []*regexp.Regexp, body string) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(body) {
			return true
		}
	}
	return false
}

func (f *MessageFilter) matchJSONFilter(body []byte) bool {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false
	}

	iter := f.jsonFilter.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

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

// GetCompilationErrors returns any errors encountered during pattern compilation
func (f *MessageFilter) GetCompilationErrors() []error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.compilationErrors
}
