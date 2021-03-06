package main

import (
	"fmt"
	"os"

	"github.com/Nivl/voorhees/internal/voorhees"
)

// Version contains the current version of the app
var version = "DEV"

func main() {
	flags, err := voorhees.ParseFlags(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse the flags: %s\n", err.Error())
		os.Exit(voorhees.ExitFailure)
	}

	if flags.PrintHelp {
		fmt.Fprintln(os.Stdout, "Usage:")
		fmt.Fprint(os.Stdout, flags.Set.FlagUsages())
		os.Exit(voorhees.ExitSuccess)
	}

	if flags.PrintVersion {
		fmt.Fprintln(os.Stdout, version)
		os.Exit(voorhees.ExitSuccess)
	}

	cfg, err := voorhees.LoadConfigFromFlags(flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(voorhees.ExitFailure)
	}

	os.Exit(voorhees.Run(cfg, os.Stdin, os.Stderr))
}
