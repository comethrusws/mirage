package scenario

import (
	"mirage/internal/config"
	"net/http"
	"time"
)

// ServeMock writes the mock response to the writer
func ServeMock(w http.ResponseWriter, s *config.Scenario) {
	// Delay if configured
	if s.Response.Delay > 0 {
		time.Sleep(s.Response.Delay)
	}

	// Set Headers
	for k, v := range s.Response.Headers {
		w.Header().Set(k, v)
	}

	// Set Status defaults to 200 if 0
	status := s.Response.Status
	if status == 0 {
		status = 200
	}
	w.WriteHeader(status)

	// Write Body
	if s.Response.Body != "" {
		w.Write([]byte(s.Response.Body))
	}
}
