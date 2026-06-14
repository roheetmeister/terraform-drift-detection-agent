package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/detector"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/providers"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/reporter"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/state"
	"github.com/roheetmeister/terraform-drift-detection-agent/pkg/config"
)

// ReportStore is a function that stores a completed scan report (used by the dashboard).
type ReportStore func(r *detector.ScanReport)

// Run starts a blocking cron-scheduled drift scanner.
// It runs until the context is cancelled.
func Run(ctx context.Context, cfg *config.Config, p providers.Provider, store ReportStore) error {
	c := cron.New()

	_, err := c.AddFunc(cfg.CronExpr, func() {
		r, err := runScan(ctx, cfg, p)
		if err != nil {
			log.Printf("scan error: %v", err)
			return
		}
		reporter.PrintSummary(r)
		if store != nil {
			store(r)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid cron expression %q: %w", cfg.CronExpr, err)
	}

	fmt.Printf("Scheduler started with cron: %s\n", cfg.CronExpr)
	fmt.Println("Press Ctrl+C to stop.")

	c.Start()
	<-ctx.Done()
	c.Stop()
	return nil
}

func runScan(ctx context.Context, cfg *config.Config, p providers.Provider) (*detector.ScanReport, error) {
	st, err := state.Parse(ctx, cfg.StatePath)
	if err != nil {
		return nil, fmt.Errorf("parsing state: %w", err)
	}

	report, err := detector.Run(ctx, st, p, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("running detector: %w", err)
	}
	report.StateFile = cfg.StatePath
	return report, nil
}
