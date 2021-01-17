package main

import (
	"os"

	"github.com/Nivl/voorhees/internal/voorhees"
)

func main() {
	os.Exit(voorhees.Run(os.Args, os.Stdin, os.Stderr))
}
