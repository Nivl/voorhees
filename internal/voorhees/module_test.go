package voorhees_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nivl/voorhees/internal/voorhees"
	"github.com/stretchr/testify/require"
)

func TestParseJSON(t *testing.T) {
	t.Parallel()

	t.Run("valid json", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "go-list-output.txt"))
		require.NoError(t, err, "os.Open() was expected to succeed")
		t.Cleanup(func() {
			require.NoError(t, f.Close())
		})

		modules, err := voorhees.ParseGoList(f)
		require.NoError(t, err, "ParseGoList() was expected to succeed")
		require.Len(t, modules, 11)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "go-list-output.txt"))
		require.NoError(t, err, "os.Open() was expected to succeed")
		t.Cleanup(func() {
			require.NoError(t, f.Close())
		})

		// we remove the first char, making the json invalid
		_, err = f.Seek(1, io.SeekStart)
		require.NoError(t, err, "Seek() was expected to succeed")

		_, err = voorhees.ParseGoList(f)
		require.Error(t, err, "ParseJSON() was expected to fail")
	})
}
