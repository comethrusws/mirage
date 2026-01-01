package scenario

import (
	"mirage/internal/config"
	"net/http"
	"path/filepath"
)

// Matcher checks if a request matches a scenario
type Matcher struct {
	Scenarios []config.Scenario
}

// NewMatcher creates a new matcher with the given scenarios
func NewMatcher(scenarios []config.Scenario) *Matcher {
	return &Matcher{Scenarios: scenarios}
}

// Match finds the first matching scenario for the request
func (m *Matcher) Match(r *http.Request) *config.Scenario {
	for i := range m.Scenarios {
		s := &m.Scenarios[i]
		if matches(s, r) {
			return s
		}
	}
	return nil
}

func matches(s *config.Scenario, r *http.Request) bool {
	// Match Method
	if s.Match.Method != "" && s.Match.Method != r.Method {
		return false
	}

	// Match Path
	// Handle wildcard paths
	reqPath := r.URL.Path
	if s.Match.Path != "" {
		// Use filepath.Match for glob support (Note: logic might vary for URL paths vs file paths but simple globs * work)
		// Or simple prefix/suffix check. The requirement says "match by path".
		// Example: /api/*
		matched, _ := filepath.Match(s.Match.Path, reqPath)
		if !matched && s.Match.Path != reqPath {
			return false
		}
	}

	// Match Headers
	for k, v := range s.Match.Headers {
		if r.Header.Get(k) != v {
			return false
		}
	}

	return true
}
