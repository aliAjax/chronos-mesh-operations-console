package control

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chronos-mesh/chronos-mesh/internal/clockmodel"
	"github.com/chronos-mesh/chronos-mesh/internal/leapsecond"
	"github.com/chronos-mesh/chronos-mesh/internal/ratelimit"
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"net/http"
	"time"
)

type Server struct {
	Addr    string
	Model   *clockmodel.Model
	Sources *source.Manager
	Leap    *leapsecond.Table
	Limit   *ratelimit.Limiter
	Started time.Time
}

func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.health)
	m.HandleFunc("/readyz", s.ready)
	m.HandleFunc("/metrics", s.metrics)
	m.HandleFunc("/v1/time", s.time)
	m.HandleFunc("/v1/sources", s.sources)
	m.HandleFunc("/v1/leap", s.leap)
	return m
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "uptime": time.Since(s.Started).String()})
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	st := s.Model.State()
	ready := st.Synchronized && s.Leap.Valid(time.Now())
	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(map[string]any{"synchronized": st.Synchronized, "stratum": st.Stratum, "leap_valid": s.Leap.Valid(time.Now()), "uncertainty": st.Uncertainty.String()})
}
func (s *Server) time(w http.ResponseWriter, r *http.Request) {
	x := s.Model.Snapshot()
	json.NewEncoder(w).Encode(map[string]any{"now": x.Now, "stratum": x.Stratum, "reference": x.Reference, "synchronized": x.Synchronized, "root_delay": x.RootDelay.String(), "root_dispersion": x.RootDispersion.String()})
}
func (s *Server) sources(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.Sources.List())
}
func (s *Server) leap(w http.ResponseWriter, r *http.Request) {
	v, e, d, entries := s.Leap.Snapshot()
	json.NewEncoder(w).Encode(map[string]any{"version": v, "expires": e, "digest": d, "entries": entries})
}
func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	st := s.Model.State()
	fmt.Fprintf(w, "# TYPE chronos_synchronized gauge\nchronos_synchronized %d\nchronos_stratum %d\nchronos_uncertainty_seconds %f\n", boolInt(st.Synchronized), st.Stratum, st.Uncertainty.Seconds())
}
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
func Run(ctx context.Context, s *Server) error {
	h := &http.Server{Addr: s.Addr, Handler: s.Handler(), ReadHeaderTimeout: 3 * time.Second, IdleTimeout: 30 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		h.Shutdown(shutdownCtx)
	}()
	return h.ListenAndServe()
}
