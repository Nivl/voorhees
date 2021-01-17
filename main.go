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

	if flags.PrintVersion {
		fmt.Fprintln(os.Stdout, version)
		os.Exit(voorhees.ExitSuccess)
	}

	os.Exit(voorhees.Run(flags, os.Stdin, os.Stderr))
}
