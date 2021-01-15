package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Nivl/voorhees/internal/modutil"
	"github.com/olekukonko/tablewriter"
)

// Results contains all the modules that need to be reported
type Results struct {
	Unmaintained []*modutil.Module
}

// HasModules checks if the results contains any modules
func (r *Results) HasModules() bool {
	return len(r.Unmaintained) > 0
}

func (r *Results) print(w io.Writer) {
	needSpacing := false
	if len(r.Unmaintained) > 0 {
		if needSpacing {
			fmt.Fprintln(w)
		}

		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Module", "Last Update", "Indirect"})
		table.SetColumnAlignment([]int{
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_CENTER,
			tablewriter.ALIGN_CENTER,
		})
		for _, m := range r.Unmaintained {
			lastUpdateTime := m.Time
			if m.HasUpdate() {
				lastUpdateTime = m.Update.Time
			}
			monthsPassed := time.Since(*lastUpdateTime) / (24 * time.Hour) / 30
			table.Append([]string{
				m.Path,
				fmt.Sprintf("%d months ago (%s)", monthsPassed, m.Time.Format("2006/01/02")),
				strconv.FormatBool(m.Indirect),
			})
		}
		table.Render()
	}
}
