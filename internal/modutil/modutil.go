// Package modutil contains various struct and functions to work on mod files
package modutil

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Module represents a single Go module
// Copied from `go help list`:
// https://github.com/golang/go/blob/e5f0d144f96c24f9244590a5414c402a10a1aba0/src/cmd/go/internal/list/list.go#L204
type Module struct {
	Path     string       // module path
	Version  string       // module version
	Versions []string     // available module versions (with -versions)
	Replace  *Module      // replaced by this module
	Time     *time.Time   // time version was created
	Update   *Module      // available update, if any (with -u)
	Main     bool         // is this the main module?
	Indirect bool         // is this module only an indirect dependency of main module?
	Dir      string       // directory holding files for this module, if any
	GoMod    string       // path to go.mod file for this module, if any
	Error    *ModuleError // error loading module
}

// HasUpdate Returns whether the module has an update available or not
func (m *Module) HasUpdate() bool {
	if m.Update == nil || m.Update.Time == nil {
		return false
	}
	if m.Time == nil {
		return true
	}
	// It's possible that a tag appears as an "update" from a commit, even
	// if that tag is older
	return m.Update.Time.After(*m.Time)
}

// ModuleError contains the error message that occurred when loading the module
type ModuleError struct {
	Err string
}

// ParseCwd parses the mod file from the current working directory
func ParseCwd() ([]*Module, error) {
	// get an invalid JSON list of all modules
	// The command returns one json encoded module per line
	out, err := run("go", "list", "-m", "-u", "-json", "all")
	if err != nil {
		return nil, err
	}
	return ParseJSON(out)
}

// ParseJSON parses the JSON output of `go list` and returns a list of module
func ParseJSON(golistOutput string) ([]*Module, error) {
	// make list a valid JSON list
	golistOutput = "[" + golistOutput + "]"
	golistOutput = strings.ReplaceAll(golistOutput, "}\n{", "},\n{")

	modules := []*Module{}
	err := json.Unmarshal([]byte(golistOutput), &modules)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the JSON output of go list")
	}
	return modules, nil
}
