package scenario

import (
	"mirage/internal/config"
	"net/http"
	"path/filepath"
)

type RuntimeScenario struct {
	config.Scenario
	Enabled bool
}

type Matcher struct {
	Scenarios []*RuntimeScenario
}

func NewMatcher(scenarios []config.Scenario) *Matcher {
	runtimeScenarios := make([]*RuntimeScenario, len(scenarios))
	for i, s := range scenarios {
		runtimeScenarios[i] = &RuntimeScenario{
			Scenario: s,
			Enabled:  true,
		}
	}
	return &Matcher{Scenarios: runtimeScenarios}
}

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

func (m *Matcher) SetEnabled(name string, enabled bool) bool {
	for _, s := range m.Scenarios {
		if s.Name == name {
			s.Enabled = enabled
			return true
		}
	}
	return false
}

func (m *Matcher) GetScenarios() []RuntimeScenario {
	res := make([]RuntimeScenario, len(m.Scenarios))
	for i, s := range m.Scenarios {
		res[i] = *s
	}
	return res
}

func matches(s *config.Scenario, r *http.Request) bool {
	if s.Match.Method != "" && s.Match.Method != r.Method {
		return false
	}

	reqPath := r.URL.Path
	if s.Match.Path != "" {
		matched, _ := filepath.Match(s.Match.Path, reqPath)
		if !matched && s.Match.Path != reqPath {
			return false
		}
	}

	for k, v := range s.Match.Headers {
		if r.Header.Get(k) != v {
			return false
		}
	}

	return true
}
