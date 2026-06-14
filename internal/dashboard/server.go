package dashboard

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/detector"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/providers"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/state"
	"github.com/roheetmeister/terraform-drift-detection-agent/pkg/config"
)

//go:embed static
var staticFiles embed.FS

// Server holds state for the dashboard HTTP server.
type Server struct {
	cfg      *config.Config
	provider providers.Provider
	mu       sync.RWMutex
	reports  []*detector.ScanReport
}

// New creates a new dashboard server.
func New(cfg *config.Config, p providers.Provider) *Server {
	return &Server{cfg: cfg, provider: p}
}

// AddReport stores a scan report (called by the scheduler).
func (s *Server) AddReport(r *detector.ScanReport) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reports = append(s.reports, r)
	if len(s.reports) > s.cfg.MaxReports {
		s.reports = s.reports[len(s.reports)-s.cfg.MaxReports:]
	}
}

// Start runs the HTTP server until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Serve embedded static files at /
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("embedding static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API: list all scan reports
	mux.HandleFunc("/api/reports", s.handleReports)

	// API: trigger an on-demand scan
	mux.HandleFunc("/api/scan", s.handleScan)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutCtx)
	}()

	fmt.Printf("Dashboard running at http://localhost:%d\n", s.cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("dashboard server: %w", err)
	}
	return nil
}

func (s *Server) handleReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeJSON(w, s.reports)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.cfg.StatePath == "" {
		http.Error(w, "no state file configured (start drift serve --state <path>)", http.StatusBadRequest)
		return
	}

	st, err := state.Parse(r.Context(), s.cfg.StatePath)
	if err != nil {
		log.Printf("scan error (state parse): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	report, err := detector.Run(r.Context(), st, s.provider, s.cfg.Region)
	if err != nil {
		log.Printf("scan error (detector): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	report.StateFile = s.cfg.StatePath

	s.AddReport(report)
	writeJSON(w, report)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}
