package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getExecModifyTime(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	modtime, err := getExecModifyTime()
	require.NoError(err)
	require.NotZero(modtime)
}
