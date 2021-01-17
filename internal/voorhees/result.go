package voorhees

import (
	"fmt"
	"io"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Results contains all the modules that need to be reported
type Results struct {
	baseTime     time.Time
	Unmaintained []*Module
}

// NewResults returns a Results object ready to be filled
func NewResults() *Results {
	return &Results{
		baseTime: time.Now(),
	}
}

// HasModules checks if the results contains any modules
func (r *Results) HasModules() bool {
	return len(r.Unmaintained) > 0
}

func (r *Results) print(w io.Writer) {
	if len(r.Unmaintained) > 0 {
		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Module", "Last Update"})
		table.SetColumnAlignment([]int{
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_CENTER,
		})
		for _, m := range r.Unmaintained {
			lastUpdateTime := m.Time
			if m.HasUpdate() {
				lastUpdateTime = m.Update.Time
			}
			monthsPassed := r.baseTime.Sub(*lastUpdateTime) / (24 * time.Hour) / 30
			table.Append([]string{
				m.Path,
				fmt.Sprintf("%d months ago (%s)", monthsPassed, lastUpdateTime.Format("2006/01/02")),
			})
		}
		table.Render()
	}
}
