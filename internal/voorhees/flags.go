package voorhees

import flag "github.com/cornfeedhobo/pflag"

// Flags represents all the flags accepted by the CLI
type Flags struct {
	IgnoredPkgs    []string
	MaxMonths      int
	PrintVersion   bool
	PrintHelp      bool
	ConfigFilePath string

	Set *flag.FlagSet
}

// ParseFlags parses the provided arguments (os.Args) and extracts the flags
func ParseFlags(args []string) (*Flags, error) {
	flags := &Flags{}
	flags.Set = flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Set.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "Coma separated list of packages to ignore")
	flags.Set.IntVarP(&flags.MaxMonths, "limit", "l", 6, "Number of months after which a dep is considered unmaintained")
	flags.Set.BoolVarP(&flags.PrintVersion, "version", "v", false, "Print version")
	flags.Set.BoolVarP(&flags.PrintHelp, "help", "h", false, "Print help")
	flags.Set.StringVarP(&flags.ConfigFilePath, "config-file", "c", DefaultConfigFilePath, "path to the config file")
	return flags, flags.Set.Parse(args)
}
