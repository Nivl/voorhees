package main

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/Nivl/voorhees/internal/modutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description    string
		argv           []string
		expectedResult Flags
		expectedError  error
	}{
		{
			description: "default flags",
			argv:        []string{"bin"},
			expectedResult: Flags{
				MaxWeeks:    26,
				IgnoredPkgs: []string{},
			},
			expectedError: nil,
		},
		{
			description: "set all",
			argv: []string{
				"bin",
				"-l",
				"4",
				"-i",
				"pkg1,pkg2",
				"--ignore=pkg3,pkg4",
			},
			expectedResult: Flags{
				MaxWeeks:    4,
				IgnoredPkgs: []string{"pkg1", "pkg2", "pkg3", "pkg4"},
			},
			expectedError: nil,
		},
		{
			description: "invalid flag",
			argv: []string{
				"bin",
				"--nope",
			},
			expectedResult: Flags{},
			expectedError:  errors.New("unknown flag: --nope"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			flags, err := parseFlags(tc.argv)
			if tc.expectedError != nil {
				require.Error(t, err, "parseFlags should have failed")
				require.Equal(t, tc.expectedError, err, "parseFlags failed with an unexpected error")
				return
			}

			require.NoError(t, err, "parseFlags should have succeed")
			assert.Equal(t, tc.expectedResult, *flags)
		})
	}
}

func TestParseModules(t *testing.T) {
	t.Parallel()

	now := time.Now()
	TwoMonthsAgo := now.Add(-8 * 7 * 24 * time.Hour)
	OneYearAgo := now.Add(-366 * 24 * time.Hour)

	validModule := &modutil.Module{Path: "valid/pkg", Version: "0.0.1", Time: &TwoMonthsAgo}
	indirectModule := &modutil.Module{Path: "indirect/pkg", Indirect: true, Time: &TwoMonthsAgo}
	updatedToModule := &modutil.Module{Path: "updated/pkg", Version: "1.0.0", Time: &now}
	updatedModule := &modutil.Module{Path: "updated/pkg", Version: "0.0.1", Time: &OneYearAgo, Update: updatedToModule}
	oldModule := &modutil.Module{Path: "old/pkg", Version: "0.0.1", Time: &OneYearAgo}

	testCases := []struct {
		description string
		flags       Flags
		modules     []*modutil.Module
		expected    Results
	}{
		{description: "no modules"},
		{
			description: "No updates in the last 6 months",
			flags: Flags{
				MaxWeeks: 26,
			},
			modules: []*modutil.Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*modutil.Module{
					oldModule,
				},
			},
		},
		{
			description: "No updates in the last month",
			flags: Flags{
				MaxWeeks: 4,
			},
			modules: []*modutil.Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*modutil.Module{
					validModule,
					oldModule,
				},
			},
		},
		{
			description: "No updates in the last month with ignore package",
			flags: Flags{
				MaxWeeks:    4,
				IgnoredPkgs: []string{validModule.Path},
			},
			modules: []*modutil.Module{
				validModule,
				indirectModule,
				updatedModule,
				oldModule,
			},
			expected: Results{
				Unmaintained: []*modutil.Module{
					oldModule,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			week := 7 * 24 * time.Hour
			expirationDate := time.Now().Add(-time.Duration(tc.flags.MaxWeeks) * week)

			res := parseModules(&tc.flags, expirationDate, tc.modules)
			assert.Equal(t, tc.expected, *res)
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		args         []string
		expectedBuf  string
		expectedCode int
	}{
		{
			description:  "invalid flags",
			args:         []string{"bin", "--nope"},
			expectedCode: ExitFailure,
			expectedBuf:  "could not parse the flags: unknown flag: --nope\n",
		},
		{
			description:  "limit too low",
			args:         []string{"bin", "-l", "2"},
			expectedCode: ExitFailure,
			expectedBuf:  "the limit cannot be below 4\n",
		},
		{
			description: "happy path",
			// We're cheating a bit by ignoring all the packages
			args:         []string{"bin", "-i", "github.com"},
			expectedCode: ExitSuccess,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			buf := bytes.Buffer{}
			w := bufio.NewWriter(&buf)
			exitStatus := run(tc.args, w)
			require.NoError(t, w.Flush(), "Flush() should have work")
			assert.Equal(t, tc.expectedCode, exitStatus)
			if tc.expectedBuf != "" {
				assert.Equal(t, tc.expectedBuf, buf.String())
			}
		})
	}
}
