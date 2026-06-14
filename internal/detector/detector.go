package detector

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/providers"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/state"
)

// Run performs a full drift scan against the given state, using the provider to fetch live data.
func Run(ctx context.Context, st *state.TFState, p providers.Provider, region string) (*ScanReport, error) {
	report := &ScanReport{
		ScannedAt: time.Now(),
		Provider:  p.Name(),
		Region:    region,
	}

	supported := make(map[string]bool)
	for _, t := range p.SupportedTypes() {
		supported[t] = true
	}

	for _, res := range st.Resources {
		if res.Mode != "managed" {
			continue
		}
		if !supported[res.Type] {
			continue
		}
		for _, inst := range res.Instances {
			expected := normalizer.FromState(&res, &inst, region)
			report.Summary.Total++

			actual, err := p.FetchResource(ctx, res.Type, expected.ID, region)
			if err != nil {
				// Treat fetch errors as missing (resource may have been deleted)
				report.Drifts = append(report.Drifts, DriftResult{
					ResourceID:   expected.ID,
					ResourceType: res.Type,
					ResourceName: res.Name,
					DriftType:    DriftMissing,
					ScannedAt:    time.Now(),
				})
				report.Summary.Missing++
				report.Summary.Drifted++
				continue
			}

			dr := compare(expected, actual, res.Name)
			if dr != nil {
				report.Drifts = append(report.Drifts, *dr)
				report.Summary.Drifted++
				switch dr.DriftType {
				case DriftMissing:
					report.Summary.Missing++
				case DriftModified:
					report.Summary.Modified++
				case DriftTagged:
					report.Summary.Tagged++
				}
			} else {
				report.Summary.Clean++
			}
		}
	}

	return report, nil
}

// compare returns a DriftResult if expected and actual differ, or nil if clean.
func compare(expected, actual *normalizer.Resource, resourceName string) *DriftResult {
	if actual == nil {
		return &DriftResult{
			ResourceID:   expected.ID,
			ResourceType: expected.Type,
			ResourceName: resourceName,
			DriftType:    DriftMissing,
			ScannedAt:    time.Now(),
		}
	}

	var diffs []AttributeDiff

	// Compare key scalar attributes (exclude computed/noisy fields)
	skip := map[string]bool{
		"id": true, "arn": true, "tags_all": true,
		"timeouts": true, "tags": true,
	}
	for k, ev := range expected.Attributes {
		if skip[k] {
			continue
		}
		av, exists := actual.Attributes[k]
		if !exists {
			continue // live API may not return every TF attribute
		}
		es := fmt.Sprintf("%v", ev)
		as := fmt.Sprintf("%v", av)
		if es != as {
			diffs = append(diffs, AttributeDiff{Attribute: k, Expected: ev, Actual: av})
		}
	}

	// Compare tags: state → cloud (missing or changed values)
	for k, ev := range expected.Tags {
		av, exists := actual.Tags[k]
		if !exists {
			diffs = append(diffs, AttributeDiff{Attribute: "tag:" + k, Expected: ev, Actual: nil})
		} else if ev != av {
			diffs = append(diffs, AttributeDiff{Attribute: "tag:" + k, Expected: ev, Actual: av})
		}
	}

	// Compare tags: cloud → state (manually added tags not in state)
	for k, av := range actual.Tags {
		if _, exists := expected.Tags[k]; !exists {
			diffs = append(diffs, AttributeDiff{Attribute: "tag:" + k, Expected: nil, Actual: av})
		}
	}

	if len(diffs) == 0 {
		return nil
	}

	driftType := DriftModified
	allTags := true
	for _, d := range diffs {
		if !strings.HasPrefix(d.Attribute, "tag:") {
			allTags = false
			break
		}
	}
	if allTags {
		driftType = DriftTagged
	}

	return &DriftResult{
		ResourceID:   expected.ID,
		ResourceType: expected.Type,
		ResourceName: resourceName,
		DriftType:    driftType,
		Differences:  diffs,
		ScannedAt:    time.Now(),
	}
}
