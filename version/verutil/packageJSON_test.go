package verutil

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersionFromSource(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	projectRootDir := filepath.Join("..", "..")
	packageFilePath := filepath.Join("test", "package.json")
	version, err := GetVersionFromSource(projectRootDir, packageFilePath)
	require.NoError(err)
	require.EqualValues("0.0.1", version.Version)
	require.NotEmpty(version.GoVersion)
	require.NotEmpty(version.BuildTime)
	require.NotEmpty(version.GitCommit)
}
