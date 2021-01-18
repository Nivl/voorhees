package voorhees_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nivl/voorhees/internal/voorhees"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "config", "1", "valid.yml"))
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, f.Close())
		})
		cfg, err := voorhees.NewConfig(f)
		require.NoError(t, err)

		assert.True(t, cfg.IsIgnored("pkg.tld/skipped"), "package should be ignored")
		assert.True(t, cfg.IsIgnored("pkg.tld/ignored"), "package should be ignored")

		tenMonths := 10 * 30 * 24 * time.Hour
		assert.Equal(t, tenMonths, cfg.Duration("pkg.tld/months"), "package should be >= 0")

		FiftyTwoWeeks := 52 * 7 * 24 * time.Hour
		assert.Equal(t, FiftyTwoWeeks, cfg.Duration("pkg.tld/weeks"), "package should be >= 0")
	})

	t.Run("failure cases", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			desc          string
			filePath      string
			expectedError string
		}{
			{
				desc:          "should fail on invalid keys",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_keys.yml"),
				expectedError: "doesNotExists not found",
			},
			{
				desc:          "should fail on unknown values",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value.yml"),
				expectedError: "invalid rule value: unknownValue",
			},
			{
				desc:          "should fail on invalid duration number",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value_duration_number.yml"),
				expectedError: "expected a number > 0",
			},
			{
				desc:          "should fail on invalid duration type",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value_duration_type.yml"),
				expectedError: "unexpected duration type: days",
			},
		}
		for i, tc := range testCases {
			tc := tc
			i := i
			t.Run(fmt.Sprintf("%d/%s", i, tc.desc), func(t *testing.T) {
				t.Parallel()

				f, err := os.Open(tc.filePath)
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, f.Close())
				})
				_, err = voorhees.NewConfig(f)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})
}
