package scenario

import (
	"mirage/internal/config"
	"net/http"
	"time"
)

func ServeMock(w http.ResponseWriter, s *config.Scenario) {
	if s.Response.Delay > 0 {
		time.Sleep(s.Response.Delay)
	}

	for k, v := range s.Response.Headers {
		w.Header().Set(k, v)
	}

	status := s.Response.Status
	if status == 0 {
		status = 200
	}
	w.WriteHeader(status)

	if s.Response.Body != "" {
		w.Write([]byte(s.Response.Body))
	}
}
