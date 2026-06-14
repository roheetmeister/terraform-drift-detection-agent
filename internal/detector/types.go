package detector

import "time"

// DriftType classifies the kind of drift detected.
type DriftType string

const (
	DriftMissing  DriftType = "missing"     // resource in state but not in cloud
	DriftModified DriftType = "modified"    // resource exists but attributes differ
	DriftTagged   DriftType = "tag_changed" // only tag differences
	DriftExtra    DriftType = "extra"       // resource in cloud but not in state
)

// AttributeDiff holds a single attribute-level difference.
type AttributeDiff struct {
	Attribute string      `json:"attribute"`
	Expected  interface{} `json:"expected"`
	Actual    interface{} `json:"actual"`
}

// DriftResult describes drift found for a single resource.
type DriftResult struct {
	ResourceID   string          `json:"resource_id"`
	ResourceType string          `json:"resource_type"`
	ResourceName string          `json:"resource_name"`
	DriftType    DriftType       `json:"drift_type"`
	Differences  []AttributeDiff `json:"differences,omitempty"`
	ScannedAt    time.Time       `json:"scanned_at"`
}

// ScanSummary aggregates counts from a scan run.
type ScanSummary struct {
	Total    int `json:"total_resources"`
	Drifted  int `json:"drifted"`
	Missing  int `json:"missing"`
	Modified int `json:"modified"`
	Tagged   int `json:"tag_changed"`
	Clean    int `json:"clean"`
}

// ScanReport is the full output of one drift scan.
type ScanReport struct {
	StateFile string       `json:"state_file"`
	Provider  string       `json:"provider"`
	Region    string       `json:"region"`
	ScannedAt time.Time    `json:"scanned_at"`
	Summary   ScanSummary  `json:"summary"`
	Drifts    []DriftResult `json:"drifts"`
}
