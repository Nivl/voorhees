package voorhees

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseModules(t *testing.T) {
	t.Parallel()

	now := time.Now()
	TwoMonthsAgo := now.Add(-8 * 7 * 24 * time.Hour)
	OneYearAgo := now.Add(-366 * 24 * time.Hour)

	validModule := &Module{Path: "valid/pkg", Version: "0.0.1", Time: &TwoMonthsAgo}
	indirectModule := &Module{Path: "indirect/pkg", Indirect: true, Time: &TwoMonthsAgo}
	updatedToModule := &Module{Path: "updated/pkg", Version: "1.0.0", Time: &now}
	updatedModule := &Module{Path: "updated/pkg", Version: "0.0.1", Time: &OneYearAgo, Update: updatedToModule}
	oldModule := &Module{Path: "old/pkg", Version: "0.0.1", Time: &OneYearAgo}

	testCases := []struct {
		description string
		cfg         *Config
		modules     []*Module
		expected    Results
	}{
		{description: "no modules"},
		{
			description: "No updates in the last 6 months",
			cfg:         newDefaultConfig(),
			modules: []*Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*Module{
					oldModule,
				},
			},
		},
		{
			description: "No updates in the last month",
			cfg: &Config{
				defaultLimit: 4 * month,
				limits:       map[string]time.Duration{},
				toSkip:       map[string]struct{}{},
			},
			modules: []*Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*Module{
					validModule,
					oldModule,
				},
			},
		},
		{
			description: "No updates in the last month with ignore package",
			cfg: &Config{
				defaultLimit: 4 * month,
				limits:       map[string]time.Duration{},
				toSkip: map[string]struct{}{
					validModule.Path: {},
				},
			},
			modules: []*Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*Module{
					oldModule,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			res := parseModules(tc.cfg, tc.modules)
			assert.Equal(t, tc.expected.Unmaintained, res.Unmaintained)
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		cfg          *Config
		expectedBuf  string
		expectedCode int
	}{
		{
			description: "happy path",
			// We're cheating a bit by ignoring all the packages
			cfg: &Config{
				defaultLimit: 6 * month,
				limits:       map[string]time.Duration{},
				toSkip: map[string]struct{}{
					"github.com/olekukonko/tablewriter": {},
					"github.com/pkg/errors":             {},
					"github.com/spf13/pflag":            {},
					"github.com/stretchr/testify":       {},
				},
			},
			expectedCode: ExitSuccess,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(filepath.Join("testdata", "go-list-output.txt"))
			require.NoError(t, err, "os.Open() was expected to succeed")
			t.Cleanup(func() {
				require.NoError(t, f.Close())
			})

			buf := bytes.Buffer{}
			w := bufio.NewWriter(&buf)
			exitStatus := Run(tc.cfg, f, w)
			require.NoError(t, w.Flush(), "Flush() should have work")
			assert.Equal(t, tc.expectedCode, exitStatus)
			if tc.expectedBuf != "" {
				assert.Equal(t, tc.expectedBuf, buf.String())
			}
		})
	}
}
