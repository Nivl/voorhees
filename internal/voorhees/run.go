package voorhees

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Run runs Voorhees
// - args is expected to contain os.Args
// - in is expected to contain the output of go list (stdin most likely)
// - out is expected to be were errors will be printed (stderr)
func Run(args []string, in io.Reader, out io.Writer) (exitStatus int) {
	flags, err := parseFlags(args)
	if err != nil {
		fmt.Fprintf(out, "could not parse the flags: %s\n", err.Error())
		return ExitFailure
	}

	// Doesn't make sense to treat a dep unmaintained after a couple of weeks
	if flags.MaxWeeks < 4 {
		fmt.Fprintln(out, "the limit cannot be below 4")
		return ExitFailure
	}
	week := 7 * 24 * time.Hour
	expirationDate := time.Now().Add(-time.Duration(flags.MaxWeeks) * week)

	modules, err := parseGoList(os.Stdin)
	if err != nil {
		fmt.Fprintf(out, "could not parse the go.mod file: %s\n", err.Error())
		return ExitFailure
	}

	res := parseModules(flags, expirationDate, modules)
	res.print(out)

	if res.HasModules() {
		return ExitFailure
	}
	return ExitSuccess
}

func parseModules(f *Flags, expirationDate time.Time, modules []*Module) *Results {
	res := NewResults()
	for _, m := range modules {
		// skip ignored packages
		isIgnored := false
		for _, pkg := range f.IgnoredPkgs {
			if strings.HasPrefix(m.Path, pkg) {
				isIgnored = true
				break
			}
		}
		if isIgnored {
			continue
		}

		// skip indirects since we can't really parse this properly
		if m.Indirect {
			continue
		}

		// Report if the package hasn't been updated since the last
		// X weeks
		if m.Time != nil && m.Time.Before(expirationDate) {
			if !m.HasUpdate() || m.Update.Time.Before(expirationDate) {
				res.Unmaintained = append(res.Unmaintained, m)
				continue
			}
		}
	}

	return res
}
