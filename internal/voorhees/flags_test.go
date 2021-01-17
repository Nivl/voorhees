package voorhees_test

import (
	"errors"
	"testing"

	"github.com/Nivl/voorhees/internal/voorhees"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description    string
		argv           []string
		expectedResult voorhees.Flags
		expectedError  error
	}{
		{
			description: "default flags",
			argv:        []string{"bin"},
			expectedResult: voorhees.Flags{
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
			expectedResult: voorhees.Flags{
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
			expectedResult: voorhees.Flags{},
			expectedError:  errors.New("unknown flag: --nope"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			flags, err := voorhees.ParseFlags(tc.argv)
			if tc.expectedError != nil {
				require.Error(t, err, "ParseFlags should have failed")
				require.Equal(t, tc.expectedError, err, "ParseFlags failed with an unexpected error")
				return
			}

			require.NoError(t, err, "ParseFlags should have succeed")
			assert.NotNil(t, flags.Set)
			// We're going to cheat a bit here because we don't care
			// about this value (beside that it's not nil)
			tc.expectedResult.Set = flags.Set
			assert.Equal(t, tc.expectedResult, *flags)
		})
	}
}
