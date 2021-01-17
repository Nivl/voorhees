package voorhees

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGoList(t *testing.T) {
	t.Parallel()

	t.Run("valid json", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "go-list-output.txt"))
		require.NoError(t, err, "os.Open() was expected to succeed")
		t.Cleanup(func() {
			require.NoError(t, f.Close())
		})

		modules, err := parseGoList(f)
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

		_, err = parseGoList(f)
		require.Error(t, err, "ParseGoList() was expected to fail")
	})
}
