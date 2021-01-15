// Package modutil contains various struct and functions to work on mod files
package modutil_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/Nivl/voorhees/internal/modutil"
	"github.com/stretchr/testify/require"
)

func TestParseJSON(t *testing.T) {
	t.Parallel()

	t.Run("valid json", func(t *testing.T) {
		t.Parallel()

		content, err := ioutil.ReadFile(filepath.Join("testdata", "go-list-output.txt"))
		require.NoError(t, err, "ioutil.ReadFile() was expected to succeed")
		modules, err := modutil.ParseJSON(string(content))
		require.NoError(t, err, "ParseJSON() was expected to succeed")
		require.Len(t, modules, 96)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		content, err := ioutil.ReadFile(filepath.Join("testdata", "go-list-output.txt"))
		require.NoError(t, err, "ioutil.ReadFile() was expected to succeed")
		// we remove the first char, making the json invalid
		_, err = modutil.ParseJSON(string(content[1:]))
		require.Error(t, err, "ParseJSON() was expected to fail")
	})
}

func TestParseCwd(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		modules, err := modutil.ParseCwd()
		require.NoError(t, err, "ParseCwd() was expected to succeed")
		require.True(t, len(modules) > 0, "several modules should have been found")
	})
}
