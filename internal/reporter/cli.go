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
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	white  = color.New(color.FgWhite, color.Bold).SprintFunc()
)

// PrintCLI writes a human-readable drift report to stdout using tables throughout.
func PrintCLI(report *detector.ScanReport) {
	fmt.Printf("\n%s\n\n", cyan("=== Terraform Drift Detection Report ==="))

	// ── Scan metadata table ───────────────────────────────────────────────────
	meta := tablewriter.NewWriter(os.Stdout)
	meta.SetHeader([]string{"Field", "Value"})
	meta.SetBorder(true)
	meta.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	meta.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	meta.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{},
	)
	meta.Append([]string{"State File", report.StateFile})
	meta.Append([]string{"Provider", report.Provider})
	meta.Append([]string{"Region", report.Region})
	meta.Append([]string{"Scanned At", report.ScannedAt.Format("2006-01-02 15:04:05 UTC")})
	meta.Render()
	fmt.Println()

	// ── Summary table ─────────────────────────────────────────────────────────
	summary := tablewriter.NewWriter(os.Stdout)
	summary.SetHeader([]string{"Total", "Drifted", "Missing", "Modified", "Tag Changed", "Clean"})
	summary.SetBorder(true)
	summary.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
	)
	summary.SetColumnAlignment([]int{
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
	})
	summary.Append([]string{
		white(fmt.Sprintf("%d", report.Summary.Total)),
		colorCount(report.Summary.Drifted, red),
		colorCount(report.Summary.Missing, red),
		colorCount(report.Summary.Modified, yellow),
		colorCount(report.Summary.Tagged, yellow),
		colorCount(report.Summary.Clean, green),
	})
	summary.Render()
	fmt.Println()

	// ── Drift results table ───────────────────────────────────────────────────
	if len(report.Drifts) == 0 {
		fmt.Println(green("✔  No drift detected. Infrastructure matches state."))
		fmt.Println()
		return
	}

	fmt.Printf("%s\n\n", red(fmt.Sprintf("✘  %d resource(s) have drifted:", len(report.Drifts))))

	results := tablewriter.NewWriter(os.Stdout)
	results.SetHeader([]string{"#", "Resource ID", "Type", "Name", "Drift Type", "Attribute", "Expected", "Actual"})
	results.SetBorder(true)
	results.SetRowLine(true)
	results.SetAutoWrapText(false)
	results.SetColumnAlignment([]int{
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})
	results.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	rowNum := 1
	for _, dr := range report.Drifts {
		driftStr := driftLabel(dr.DriftType)
		if len(dr.Differences) == 0 {
			results.Append([]string{
				fmt.Sprintf("%d", rowNum),
				truncate(dr.ResourceID, 30),
				dr.ResourceType,
				dr.ResourceName,
				driftStr,
				"-", "-", "-",
			})
			rowNum++
			continue
		}
		for i, diff := range dr.Differences {
			num := ""
			rid := ""
			rtype := ""
			rname := ""
			drift := ""
			if i == 0 {
				num = fmt.Sprintf("%d", rowNum)
				rid = truncate(dr.ResourceID, 30)
				rtype = dr.ResourceType
				rname = dr.ResourceName
				drift = driftStr
				rowNum++
			}
			results.Append([]string{
				num, rid, rtype, rname, drift,
				diff.Attribute,
				truncate(fmt.Sprintf("%v", diff.Expected), 25),
				truncate(fmt.Sprintf("%v", diff.Actual), 25),
			})
		}
	}

	results.Render()
	fmt.Println()
}

func colorCount(n int, fn func(...interface{}) string) string {
	if n == 0 {
		return fmt.Sprintf("%d", n)
	}
	return fn(fmt.Sprintf("%d", n))
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
