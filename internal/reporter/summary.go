package reporter

import (
	"fmt"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/detector"
)

// PrintSummary writes a one-line summary of a scan report.
func PrintSummary(report *detector.ScanReport) {
	fmt.Printf("[%s] %s — total: %d, drifted: %d (missing: %d, modified: %d, tags: %d), clean: %d\n",
		report.ScannedAt.Format("15:04:05"),
		report.StateFile,
		report.Summary.Total,
		report.Summary.Drifted,
		report.Summary.Missing,
		report.Summary.Modified,
		report.Summary.Tagged,
		report.Summary.Clean,
	)
}
