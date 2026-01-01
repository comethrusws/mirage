package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"mirage/internal/config"
	"mirage/internal/recorder"
	"mirage/internal/scenario"
)

// Proxy implements http.Handler and forwards requests
type Proxy struct {
	client   *http.Client
	matcher  *scenario.Matcher
	recorder *recorder.Recorder
	
	// In-memory request log for dashboard
	reqLogMu   sync.RWMutex
	reqLog     []LogEntry
	MaxLogSize int
}

type LogEntry struct {
	ID        int64         `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Method    string        `json:"method"`
	URL       string        `json:"url"`
	Status    int           `json:"status"`
	Duration  time.Duration `json:"duration"`
	Matched   string        `json:"matched,omitempty"` // Scenario name if matched
}

// NewProxy creates a new Proxy instance
func NewProxy(cfg *config.Config, rec *recorder.Recorder) *Proxy {
	var m *scenario.Matcher
	if cfg != nil {
		m = scenario.NewMatcher(cfg.Scenarios)
	}

	return &Proxy{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow redirects, forward them
			},
		},
		matcher:    m,
		recorder:   rec,
		reqLog:     make([]LogEntry, 0),
		MaxLogSize: 100,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Handle Dashboard/API requests if embedded (this logic might be in separate handler, 
	// but user asked for single server on 8080. We can check prefix here or use Mux in main.
	// If main uses p for everything, p needs to dispatch.
	// BUT main uses `http.ListenAndServe(addr, p)`.
	// So p is the root handler.
	// We can check path here.
	
	// To avoid circular deps, UI handler should be passed to Proxy or handled in main via a wrapper.
	// Let's assume passed in or handled here. 
	// Easier: main creates a Mux that routes /__mirage/ to UI and / to Proxy.
	// So Proxy only handles proxy traffic.
	// main.go will change.
	
	// Proxy logic:
	
	// Read and log request body
	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Restore body
	}
	
	// Truncate body for logging if too long
	logReqBody := string(reqBody)
	if len(logReqBody) > 500 {
		logReqBody = logReqBody[:500] + "...(truncated)"
	}

	log.Printf("[REQ] %s %s Headers: %v Body: %s", r.Method, r.URL.String(), r.Header, logReqBody)
	
	var matchedScenario string
	var status int
	
	// Check for mock scenario
	if p.matcher != nil {
		if s := p.matcher.Match(r); s != nil {
			log.Printf("[MOCK] Matched scenario: %s", s.Name)
			scenario.ServeMock(w, s)
			
			duration := time.Since(start)
			log.Printf("[RES] [MOCK] Status: %d Duration: %v", s.Response.Status, duration)
			
			matchedScenario = s.Name
			status = s.Response.Status
			if status == 0 { status = 200 }
			
			p.logRequest(r, status, duration, matchedScenario)
			return
		}
	}

	outReq := r.Clone(r.Context())
	
	// Remove hop-by-hop headers
	delHopHeaders(outReq.Header)

	// Forward the request
	resp, err := p.client.Do(outReq)
	if err != nil {
		log.Printf("[ERR] Forwarding failed: %v", err)
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
		p.logRequest(r, 502, time.Since(start), "")
		return
	}
	defer resp.Body.Close()

	// Copy headers
	delHopHeaders(resp.Header)
	copyHeader(w.Header(), resp.Header)

	// Write status code
	w.WriteHeader(resp.StatusCode)
	status = resp.StatusCode

	// Read response body to log it (and write to client)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERR] Reading response body: %v", err)
		return
	}
	
	// Write body to client
	w.Write(respBody)

	// Log response
	duration := time.Since(start)
	logRespBody := string(respBody)
	if len(logRespBody) > 500 {
		logRespBody = logRespBody[:500] + "...(truncated)"
	}
	
	log.Printf("[RES] Status: %d Duration: %v Body: %s", resp.StatusCode, duration, logRespBody)

	// Record interaction if recorder is present
	if p.recorder != nil {
		p.recorder.Record(r, string(reqBody), resp, string(respBody), duration)
	}
	
	p.logRequest(r, status, duration, "")
}

func (p *Proxy) logRequest(r *http.Request, status int, duration time.Duration, matched string) {
	p.reqLogMu.Lock()
	defer p.reqLogMu.Unlock()
	
	entry := LogEntry{
		ID:        time.Now().UnixNano(),
		Timestamp: time.Now(),
		Method:    r.Method,
		URL:       r.URL.String(),
		Status:    status,
		Duration:  duration,
		Matched:   matched,
	}
	
	// prepend or append? Append is easier, loop backwards for UI
	p.reqLog = append(p.reqLog, entry)
	if len(p.reqLog) > p.MaxLogSize {
		p.reqLog = p.reqLog[1:]
	}
}

// Accessors for UI

func (p *Proxy) GetRecentRequests() []LogEntry {
	p.reqLogMu.RLock()
	defer p.reqLogMu.RUnlock()
	// Return copy
	res := make([]LogEntry, len(p.reqLog))
	copy(res, p.reqLog)
    // Reverse for display? UI can handle it.
	return res
}

func (p *Proxy) GetScenarios() []scenario.RuntimeScenario {
	if p.matcher == nil {
		return nil
	}
	// Add method in matcher to get scenarios
	return p.matcher.GetScenarios()
}

func (p *Proxy) ToggleScenario(name string, enabled bool) bool {
    if p.matcher == nil { return false }
    return p.matcher.SetEnabled(name, enabled)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	// List of hop-by-hop headers to remove
	hopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	for _, h := range hopHeaders {
		header.Del(h)
	}
}
