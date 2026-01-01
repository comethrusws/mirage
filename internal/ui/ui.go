package ui

import (
	_ "embed"
	"encoding/json"
	"net/http"
	
	"mirage/internal/proxy"
	
	"github.com/gorilla/mux"
)

//go:embed dashboard.html
var dashboardHTML string

type UI struct {
	proxy *proxy.Proxy
}

func NewUI(p *proxy.Proxy) *UI {
	return &UI{proxy: p}
}

func (u *UI) Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/__mirage/", u.handleDashboard).Methods("GET")
	r.HandleFunc("/__mirage/api/requests", u.handleRequests).Methods("GET")
	r.HandleFunc("/__mirage/api/scenarios", u.handleScenarios).Methods("GET")
	r.HandleFunc("/__mirage/api/scenarios/{name}/toggle", u.handleToggle).Methods("POST")
	return r
}

func (u *UI) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboardHTML))
}

func (u *UI) handleRequests(w http.ResponseWriter, r *http.Request) {
	logs := u.proxy.GetRecentRequests()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (u *UI) handleScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios := u.proxy.GetScenarios()
	w.Header().Set("Content-Type", "application/json")
	if scenarios == nil {
		w.Write([]byte("[]"))
		return
	}
	json.NewEncoder(w).Encode(scenarios)
}

func (u *UI) handleToggle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	success := u.proxy.ToggleScenario(name, body.Enabled)
	if !success {
		http.Error(w, "Scenario not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}
