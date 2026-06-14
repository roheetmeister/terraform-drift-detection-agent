package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/detector"
)

// PrintJSON writes the scan report as indented JSON to stdout.
func PrintJSON(report *detector.ScanReport) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("encoding report as JSON: %w", err)
	}
	return nil
}

// MarshalReport returns the report as a JSON byte slice.
func MarshalReport(report *detector.ScanReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}
