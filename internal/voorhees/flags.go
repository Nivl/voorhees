package voorhees

import flag "github.com/spf13/pflag"

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
