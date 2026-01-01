package recorder

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

// Interaction represents a captured request/response pair
type Interaction struct {
	Timestamp time.Time    `json:"timestamp"`
	Request   ReqDetail    `json:"request"`
	Response  RespDetail   `json:"response"`
	Duration  string       `json:"duration"`
}

type ReqDetail struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type RespDetail struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// Recorder captures and stores interactions
type Recorder struct {
	mu           sync.Mutex
	Interactions []Interaction
	OutputFile   string
}

// NewRecorder creates a new recorder
func NewRecorder(outputFile string) *Recorder {
	return &Recorder{
		OutputFile: outputFile,
		Interactions: make([]Interaction, 0),
	}
}

// Record captures an interaction
func (r *Recorder) Record(req *http.Request, reqBody string, resp *http.Response, respBody string, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	interaction := Interaction{
		Timestamp: time.Now(),
		Request: ReqDetail{
			Method:  req.Method,
			URL:     req.URL.String(),
			Headers: req.Header,
			Body:    reqBody,
		},
		Response: RespDetail{
			Status:  resp.StatusCode,
			Headers: resp.Header,
			Body:    respBody,
		},
		Duration: duration.String(),
	}

	r.Interactions = append(r.Interactions, interaction)
	
	// Save immediately or periodically?
	// For simplicity, save on every request for now to avoid data loss on crash
	r.save()
}

func (r *Recorder) save() error {
	data, err := json.MarshalIndent(r.Interactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.OutputFile, data, 0644)
}
