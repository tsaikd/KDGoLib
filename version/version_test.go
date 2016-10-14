package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	require.NotEmpty(String())
}

func TestJSON(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	json, err := JSON()
	require.NoError(err)
	require.NotEmpty(json)
}

func Test_getExecModifyTime(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	modtime, err := getExecModifyTime()
	require.NoError(err)
	require.NotZero(modtime)
}
