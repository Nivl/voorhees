package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Nivl/voorhees/internal/modutil"
	flag "github.com/spf13/pflag"
)

// https://www.gnu.org/software/libc/manual/html_node/Exit-Status.html
const (
	ExitSuccess = 0
	ExitFailure = 1
)

// Flags represents all the flags accepted by the CLI
type Flags struct {
	IgnoredPkgs []string
	MaxWeeks    int
}

func parseFlags(args []string) (*Flags, error) {
	flags := &Flags{}
	fset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fset.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "coma separated list of packages to ignore")
	fset.IntVarP(&flags.MaxWeeks, "limit", "l", 26, "number of weeks after which a dep is considered unmaintained")
	return flags, fset.Parse(args)
}

func main() {
	os.Exit(run(os.Args, os.Stderr))
}

func run(args []string, out io.Writer) (exitStatus int) {
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

	modules, err := modutil.ParseCwd()
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

func parseModules(f *Flags, expirationDate time.Time, modules []*modutil.Module) *Results {
	res := &Results{}
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
