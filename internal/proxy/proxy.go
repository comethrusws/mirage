package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"mirage/internal/config"
	"mirage/internal/scenario"
)

// Proxy implements http.Handler and forwards requests
type Proxy struct {
	client  *http.Client
	matcher *scenario.Matcher
}

// NewProxy creates a new Proxy instance
func NewProxy(cfg *config.Config) *Proxy {
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
		matcher: m,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

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

	// Check for mock scenario
	if p.matcher != nil {
		if s := p.matcher.Match(r); s != nil {
			log.Printf("[MOCK] Matched scenario: %s", s.Name)
			scenario.ServeMock(w, s)
			
			duration := time.Since(start)
			log.Printf("[RES] [MOCK] Status: %d Duration: %v", s.Response.Status, duration)
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
		return
	}
	defer resp.Body.Close()

	// Copy headers
	delHopHeaders(resp.Header)
	copyHeader(w.Header(), resp.Header)

	// Write status code
	w.WriteHeader(resp.StatusCode)

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
