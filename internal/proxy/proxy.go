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

type Proxy struct {
	client   *http.Client
	matcher  *scenario.Matcher
	recorder *recorder.Recorder
	
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
	Matched   string        `json:"matched,omitempty"`
}

func NewProxy(cfg *config.Config, rec *recorder.Recorder) *Proxy {
	var m *scenario.Matcher
	if cfg != nil {
		m = scenario.NewMatcher(cfg.Scenarios)
	}

	return &Proxy{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
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

	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}
	
	logReqBody := string(reqBody)
	if len(logReqBody) > 500 {
		logReqBody = logReqBody[:500] + "...(truncated)"
	}

	log.Printf("[REQ] %s %s Headers: %v Body: %s", r.Method, r.URL.String(), r.Header, logReqBody)
	
	var matchedScenario string
	var status int
	
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
	
	delHopHeaders(outReq.Header)

	resp, err := p.client.Do(outReq)
	if err != nil {
		log.Printf("[ERR] Forwarding failed: %v", err)
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
		p.logRequest(r, 502, time.Since(start), "")
		return
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)
	copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)
	status = resp.StatusCode

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERR] Reading response body: %v", err)
		return
	}
	
	w.Write(respBody)

	duration := time.Since(start)
	logRespBody := string(respBody)
	if len(logRespBody) > 500 {
		logRespBody = logRespBody[:500] + "...(truncated)"
	}
	
	log.Printf("[RES] Status: %d Duration: %v Body: %s", resp.StatusCode, duration, logRespBody)

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
	
	p.reqLog = append(p.reqLog, entry)
	if len(p.reqLog) > p.MaxLogSize {
		p.reqLog = p.reqLog[1:]
	}
}

func (p *Proxy) GetRecentRequests() []LogEntry {
	p.reqLogMu.RLock()
	defer p.reqLogMu.RUnlock()
	res := make([]LogEntry, len(p.reqLog))
	copy(res, p.reqLog)
	return res
}

func (p *Proxy) GetScenarios() []scenario.RuntimeScenario {
	if p.matcher == nil {
		return nil
	}
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
