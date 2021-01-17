package voorhees

import flag "github.com/spf13/pflag"

// Flags represents all the flags accepted by the CLI
type Flags struct {
	IgnoredPkgs  []string
	MaxWeeks     int
	PrintVersion bool
}

// ParseFlags parses the provided arguments (os.Args) and extracts the flags
func ParseFlags(args []string) (*Flags, error) {
	flags := &Flags{}
	fset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fset.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "Coma separated list of packages to ignore")
	fset.IntVarP(&flags.MaxWeeks, "limit", "l", 26, "Number of weeks after which a dep is considered unmaintained")
	fset.BoolVarP(&flags.PrintVersion, "version", "v", false, "Print version")
	return flags, fset.Parse(args)
}
