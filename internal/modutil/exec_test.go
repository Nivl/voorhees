// !build +linux
package modutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExe(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		stdout, stderr, err := exe("ls")
		require.NoError(t, err, "exe() should have succeed")
		assert.Empty(t, stderr, "stderr should have been empty")
		assert.NotEmpty(t, stdout, "stdout should have gotten data")
	})

	t.Run("commands that fails", func(t *testing.T) {
		t.Parallel()

		stdout, stderr, err := exe("ls", "/does-not-exist")
		require.Error(t, err, "exe() should have failed")
		assert.Contains(t, stderr, "No such file or directory", "unexpected stderr returned")
		assert.Empty(t, stdout, "stdout should have been empty")
	})

	t.Run("unknown command", func(t *testing.T) {
		t.Parallel()

		stdout, stderr, err := exe("kjhgfd")
		require.Error(t, err, "exe() should have failed")
		assert.Contains(t, err.Error(), "executable file not found", "unexpected error returned")
		assert.Empty(t, stderr, "stderr should have been empty")
		assert.Empty(t, stdout, "stdout should have been empty")
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		stdout, err := run("ls")
		require.NoError(t, err, "exe() should have succeed")
		assert.NotEmpty(t, stdout, "stdout should have gotten data")
	})

	t.Run("commands that fails", func(t *testing.T) {
		t.Parallel()

		stdout, err := run("ls", "/does-not-exist")
		require.Error(t, err, "exe() should have failed")
		assert.Contains(t, err.Error(), "No such file or directory", "unexpected error returned")
		assert.Empty(t, stdout, "stdout should have been empty")
	})

	t.Run("unknown command", func(t *testing.T) {
		t.Parallel()

		stdout, err := run("kjhgfd")
		require.Error(t, err, "exe() should have failed")
		assert.Contains(t, err.Error(), "executable file not found", "unexpected error returned")
		assert.Empty(t, stdout, "stdout should have been empty")
	})
}
