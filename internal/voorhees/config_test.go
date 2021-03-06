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
			expectedError error
			ErrorMatch    string
		}{
			{
				desc:       "should fail on invalid keys",
				filePath:   filepath.Join("testdata", "config", "1", "invalid_keys.yml"),
				ErrorMatch: "doesNotExists not found",
			},
			{
				desc:          "should fail on unknown values",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value.yml"),
				expectedError: voorhees.ErrConfigInvalidRuleValue,
			},
			{
				desc:          "should fail on invalid duration number",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value_duration_number.yml"),
				expectedError: voorhees.ErrConfigInvalidNumber,
			},
			{
				desc:          "should fail on invalid duration type",
				filePath:      filepath.Join("testdata", "config", "1", "invalid_value_duration_type.yml"),
				expectedError: voorhees.ErrConfigInvalidDurationType,
			},
			{
				desc:          "should fail on invalid version",
				filePath:      filepath.Join("testdata", "config", "invalid_version.yml"),
				expectedError: voorhees.ErrConfigVersion,
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
				if tc.expectedError != nil {
					assert.ErrorIs(t, err, tc.expectedError)
				}
				if tc.ErrorMatch != "" {
					assert.Contains(t, err.Error(), tc.ErrorMatch)
				}
			})
		}
	})
}

func TestLoadConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := filepath.Join("testdata", "config", "1", "valid.yml")
		cfg, err := voorhees.LoadConfigFile(p)
		require.NoError(t, err)
		assert.True(t, cfg.IsIgnored("pkg.tld/skipped"), "package should be ignored")
	})

	t.Run("should success if the default file is missing", func(t *testing.T) {
		t.Parallel()

		cfg, err := voorhees.LoadConfigFile(voorhees.DefaultConfigFilePath)
		require.NoError(t, err)
		assert.False(t, cfg.IsIgnored("pkg.tld/skipped"), "package should not be ignored")
	})

	t.Run("should fail if the file is missing", func(t *testing.T) {
		t.Parallel()

		_, err := voorhees.LoadConfigFile(".doesntexist.yml")
		require.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist, "unexpected error received")
	})
}
