package errutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConsoleErrorFormatter(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	errchild1 := New("childerror1")
	errchild2 := New("childerror2")
	errtest := New("lasterror", errchild1, errchild2)

	testDataPool := []struct {
		Formatter ErrorFormatter
		Expected  string
		Exactly   bool
		Message   string
	}{
		{
			&ConsoleFormatter{},
			"lasterror",
			true,
			"empty config",
		},
		{
			&ConsoleFormatter{
				Seperator: "; ",
			},
			"lasterror; childerror1; childerror2",
			true,
			"Seperator",
		},
		{
			&ConsoleFormatter{
				LongFile: true,
			},
			"github.com/tsaikd/KDGoLib/errutil/ConsoleFormatter_test.go lasterror",
			false,
			"LongFile",
		},
		{
			&ConsoleFormatter{
				ShortFile: true,
			},
			"ConsoleFormatter_test.go lasterror",
			true,
			"LongFile",
		},
		{
			&ConsoleFormatter{
				Seperator: "; ",
				LongFile:  true,
			},
			"github.com/tsaikd/KDGoLib/errutil/ConsoleFormatter_test.go lasterror; childerror1; childerror2",
			false,
			"Seperator|LongFile",
		},
		{
			&ConsoleFormatter{
				Seperator: "; ",
				ShortFile: true,
			},
			"ConsoleFormatter_test.go lasterror; childerror1; childerror2",
			true,
			"Seperator|ShortFile",
		},
		{
			&ConsoleFormatter{
				Seperator: "; ",
				LongFile:  true,
				Line:      true,
			},
			"github.com/tsaikd/KDGoLib/errutil/ConsoleFormatter_test.go:19 lasterror; childerror1; childerror2",
			false,
			"Seperator|LongFile|Line",
		},
		{
			&ConsoleFormatter{
				Seperator: "; ",
				ShortFile: true,
				Line:      true,
			},
			"ConsoleFormatter_test.go:19 lasterror; childerror1; childerror2",
			true,
			"Seperator|ShortFile|Line",
		},
	}

	for _, testData := range testDataPool {
		errtext, err := testData.Formatter.Format(errtest)
		require.NoError(err, testData.Message)
		if testData.Exactly {
			require.Equal(testData.Expected, errtext, testData.Message)
		} else {
			require.True(strings.HasSuffix(errtext, testData.Expected), testData.Message+" for %q", errtext)
		}
	}
}
