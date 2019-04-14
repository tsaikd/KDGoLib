package pkgutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsGoModDir(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	require.False(IsGoModDir(""))
	require.True(IsGoModDir("test"))
}

func TestParseGoMod(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	require.NoError(os.Setenv("GO111MODULE", "on"))
	result, err := ParseGoMod("test")
	require.NoError(err)
	require.Len(result, 1)
	require.EqualValues("github.com/tsaikd/KDGoLib/pkgutil/test", result[0].Path)
}
