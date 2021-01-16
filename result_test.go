package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nivl/voorhees/internal/modutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasModules(t *testing.T) {
	t.Parallel()

	// sugars
	hasModule := true

	testCases := []struct {
		description string
		res         Results
		expected    bool
	}{
		{
			description: "no modules",
			res: Results{
				Unmaintained: []*modutil.Module{},
			},
			expected: !hasModule,
		},
		{
			description: "1 unmaintained modules",
			res: Results{
				Unmaintained: []*modutil.Module{
					{},
				},
			},
			expected: hasModule,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.res.HasModules())
		})
	}
}

func TestPrint(t *testing.T) {
	t.Parallel()

	now, err := time.Parse("2006/01/02", "2019/04/30")
	require.NoError(t, err, "time.Parse() was expected to succeed")
	now = now.UTC()

	OneYearAgo := now.Add(-366 * 24 * time.Hour)
	TwoYearsAgo := now.Add(-2 * 366 * 24 * time.Hour)
	updatedToModule := &modutil.Module{Path: "updated/pkg", Version: "1.0.0", Time: &OneYearAgo}
	updatedModule := &modutil.Module{Path: "updated/pkg", Version: "0.0.1", Time: &TwoYearsAgo, Update: updatedToModule}
	oldModule := &modutil.Module{Path: "old/pkg", Version: "0.0.1", Time: &OneYearAgo}

	testCases := []struct {
		description        string
		res                Results
		expectedOutputFile string
	}{
		{
			description: "no modules",
			res: Results{
				baseTime: now,
			},
			expectedOutputFile: "empty",
		},
		{
			description:        "Unmaintained modules",
			expectedOutputFile: "test-print-unmaintained",
			res: Results{
				baseTime: now,
				Unmaintained: []*modutil.Module{
					updatedModule,
					oldModule,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			expectedOutput, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedOutputFile))
			require.NoError(t, err, "ioutil.ReadFile() was expected to succeed")

			output := bytes.Buffer{}
			w := bufio.NewWriter(&output)
			tc.res.print(w)

			require.NoError(t, w.Flush(), "Flush() should have work")

			assert.Equal(t, string(expectedOutput), output.String())
		})
	}
}
