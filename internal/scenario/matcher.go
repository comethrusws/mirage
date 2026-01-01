package scenario

import (
	"mirage/internal/config"
	"net/http"
	"path/filepath"
)

// RuntimeScenario wraps a config scenario with runtime state
type RuntimeScenario struct {
	config.Scenario
	Enabled bool
}

// Matcher checks if a request matches a scenario
type Matcher struct {
	Scenarios []*RuntimeScenario
}

// NewMatcher creates a new matcher with the given scenarios
func NewMatcher(scenarios []config.Scenario) *Matcher {
	runtimeScenarios := make([]*RuntimeScenario, len(scenarios))
	for i, s := range scenarios {
		runtimeScenarios[i] = &RuntimeScenario{
			Scenario: s,
			Enabled:  true, // Default to enabled? Or based on config? Default true.
		}
	}
	return &Matcher{Scenarios: runtimeScenarios}
}

// Match finds the first matching enabled scenario for the request
func (m *Matcher) Match(r *http.Request) *config.Scenario {
	for _, s := range m.Scenarios {
		if !s.Enabled {
			continue
		}
		if matches(&s.Scenario, r) {
			return &s.Scenario
		}
	}
	return nil
}

// ToggleScenario toggles the enabled state of a scenario
func (m *Matcher) SetEnabled(name string, enabled bool) bool {
	for _, s := range m.Scenarios {
		if s.Name == name {
			s.Enabled = enabled
			return true
		}
	}
	return false
}

// GetScenarios returns all scenarios with their state
func (m *Matcher) GetScenarios() []RuntimeScenario {
    res := make([]RuntimeScenario, len(m.Scenarios))
    for i, s := range m.Scenarios {
        res[i] = *s
    }
    return res
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
