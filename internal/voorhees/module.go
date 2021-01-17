package voorhees

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
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

// parseGoList parses the output of "go list -m -u -json all"
func parseGoList(r io.Reader) ([]*Module, error) {
	modules := []*Module{}
	d := json.NewDecoder(r)
	for i := 1; ; i++ {
		var m *Module
		if err := d.Decode(&m); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("could not parse item %d of the provided go list: %w", i, err)
		}
		modules = append(modules, m)
	}
	return modules, nil
}
