package reporter

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/detector"
)

var (
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

// PrintCLI writes a human-readable drift report to stdout.
func PrintCLI(report *detector.ScanReport) {
	fmt.Printf("\n%s\n", cyan("Terraform Drift Detection Report"))
	fmt.Printf("State file : %s\n", report.StateFile)
	fmt.Printf("Provider   : %s\n", report.Provider)
	fmt.Printf("Region     : %s\n", report.Region)
	fmt.Printf("Scanned at : %s\n\n", report.ScannedAt.Format("2006-01-02 15:04:05 UTC"))

	// Summary
	fmt.Printf("Summary  Total: %d  |  ", report.Summary.Total)
	fmt.Printf("%s  |  ", red(fmt.Sprintf("Missing: %d", report.Summary.Missing)))
	fmt.Printf("%s  |  ", yellow(fmt.Sprintf("Modified: %d", report.Summary.Modified)))
	fmt.Printf("%s  |  ", yellow(fmt.Sprintf("Tag changed: %d", report.Summary.Tagged)))
	fmt.Printf("%s\n\n", green(fmt.Sprintf("Clean: %d", report.Summary.Clean)))

	if len(report.Drifts) == 0 {
		fmt.Println(green("No drift detected. Infrastructure matches state."))
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Resource ID", "Type", "Name", "Drift", "Attribute", "Expected", "Actual"})
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetColWidth(40)

	for _, dr := range report.Drifts {
		driftStr := driftLabel(dr.DriftType)
		if len(dr.Differences) == 0 {
			table.Append([]string{
				truncate(dr.ResourceID, 25),
				dr.ResourceType,
				dr.ResourceName,
				driftStr,
				"-", "-", "-",
			})
			continue
		}
		for i, diff := range dr.Differences {
			rid := ""
			rtype := ""
			rname := ""
			drift := ""
			if i == 0 {
				rid = truncate(dr.ResourceID, 25)
				rtype = dr.ResourceType
				rname = dr.ResourceName
				drift = driftStr
			}
			table.Append([]string{
				rid, rtype, rname, drift,
				diff.Attribute,
				truncate(fmt.Sprintf("%v", diff.Expected), 30),
				truncate(fmt.Sprintf("%v", diff.Actual), 30),
			})
		}
	}

	table.Render()
	fmt.Println()
}

func driftLabel(dt detector.DriftType) string {
	switch dt {
	case detector.DriftMissing:
		return red("MISSING")
	case detector.DriftModified:
		return yellow("MODIFIED")
	case detector.DriftTagged:
		return yellow("TAG_CHANGED")
	default:
		return string(dt)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
