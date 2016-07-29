package pkgutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuessPackageFromDir(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	pkg, err := GuessPackageFromDir("")
	require.NoError(err)
	require.Contains(pkg.ImportPath, "github.com/tsaikd/KDGoLib/pkgutil")

	pkg, err = GuessPackageFromDir("/")
	require.Error(err)
	require.Nil(pkg)
}

func TestFindAllSubPackages(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	pkglist, err := FindAllSubPackages("github.com/tsaikd/KDGoLib/pkgutil", "")
	require.NoError(err)
	require.EqualValues(1, pkglist.Len())

	pkglist, err = FindAllSubPackages("github.com/tsaikd/KDGoLib", "../")
	require.NoError(err)
	require.True(pkglist.Len() > 1)
}
