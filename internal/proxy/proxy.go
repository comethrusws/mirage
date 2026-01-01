package proxy

import (
	"io"
	"log"
	"net/http"
)

// Proxy implements http.Handler and forwards requests
type Proxy struct {
	client *http.Client
}

// NewProxy creates a new Proxy instance
func NewProxy() *Proxy {
	return &Proxy{
		client: &http.Client{
			// Verify certificates for now? Or skip?
			// For a dev tool, maybe we want to be lenient, but start safe.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow redirects, forward them
			},
		},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("Received request: %s %s", r.Method, r.URL.String())

	// Handle explicit proxy requests (where URL is absolute) vs transparent/reverse (relative path)
	// For now, let's assume explicit proxy usage or we treat it as passing through.
	// If the request has no host, it might be a direct request to the proxy. Since we are a gateway,
	// we expect the client to be configured to use us as a proxy OR we are placed in front.
	// If the URL Scheme is missing, we might assume HTTP.

	outReq := r.Clone(r.Context())
	
	// If RequestURI is absolute (proxy mode), Request.URL is populated.
	// If transparent, we might need to know the target. 
	// The requirements say "Act as an HTTP/HTTPS proxy server".
	// Standard HTTP proxy receives "GET http://target.com/path HTTP/1.1".
	
	if r.URL.Scheme == "" {
		// Just a fall back for testing or transparent mode (if we knew target)
		// For now, if no scheme, we can't really forward unless we default to something?
		// But in explicit proxy mod, scheme is present.
		// Let's assume explicit proxy for this step.
		
		// If it's a CONNECT request (HTTPS), that's handled differently (tunneling).
		// Requirements: "Certificate generation/HTTPS MITM (start with HTTP only, add HTTPS later)"
		// So we only handle HTTP.
	}

	// Remove hop-by-hop headers
	delHopHeaders(outReq.Header)

	// Forward the request
	resp, err := p.client.Do(outReq)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
		log.Printf("Error forwarding: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	delHopHeaders(resp.Header)
	copyHeader(w.Header(), resp.Header)

	// Write status code
	w.WriteHeader(resp.StatusCode)

	// Copy body
	io.Copy(w, resp.Body)
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
