package recorder

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

type Interaction struct {
	Timestamp time.Time  `json:"timestamp"`
	Request   ReqDetail  `json:"request"`
	Response  RespDetail `json:"response"`
	Duration  string     `json:"duration"`
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

type Recorder struct {
	mu           sync.Mutex
	Interactions []Interaction
	OutputFile   string
}

func NewRecorder(outputFile string) *Recorder {
	return &Recorder{
		OutputFile:   outputFile,
		Interactions: make([]Interaction, 0),
	}
}

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
	r.save()
}

func (r *Recorder) save() error {
	data, err := json.MarshalIndent(r.Interactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.OutputFile, data, 0644)
}
