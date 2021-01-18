package voorhees

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Run runs Voorhees
// - in is expected to contain the output of go list (stdin most likely)
// - out is expected to be were errors will be printed (stderr)
func Run(cfg *Config, in io.Reader, out io.Writer) (exitStatus int) {
	modules, err := parseGoList(os.Stdin)
	if err != nil {
		fmt.Fprintf(out, "could not parse the go.mod file: %s\n", err.Error())
		return ExitFailure
	}

	res := parseModules(cfg, modules)
	res.print(out)

	if res.HasModules() {
		return ExitFailure
	}
	return ExitSuccess
}

func parseModules(cfg *Config, modules []*Module) *Results {
	res := NewResults()
	for _, m := range modules {
		if cfg.IsIgnored(strings.ToLower(m.Path)) {
			continue
		}

		// skip indirects since we can't really parse this properly
		if m.Indirect {
			continue
		}

		// Report if the package hasn't been updated since the last
		// X weeks
		expirationDate := time.Now().Add(-cfg.Duration(m.Path))
		if m.Time != nil && m.Time.Before(expirationDate) {
			if !m.HasUpdate() || m.Update.Time.Before(expirationDate) {
				res.Unmaintained = append(res.Unmaintained, m)
				continue
			}
		}
	}

	return res
}
