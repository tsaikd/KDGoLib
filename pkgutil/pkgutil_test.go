package pkgutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuessPackageFromDir(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	pkg, err := GuessPackageFromDir("")
	require.NoError(err)
	require.True(strings.HasSuffix(pkg.ImportPath, "github.com/tsaikd/KDGoLib/pkgutil"))

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

func TestParsePackagePaths(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	pkglist, err := ParsePackagePaths("")
	require.NoError(err)
	require.EqualValues(1, pkglist.Len())
	for pkg := range pkglist.Map() {
		require.True(strings.HasSuffix(pkg.ImportPath, "github.com/tsaikd/KDGoLib/pkgutil"))
	}

	pkglist, err = ParsePackagePaths("", "..")
	require.NoError(err)
	require.EqualValues(1, pkglist.Len())
	for pkg := range pkglist.Map() {
		require.True(strings.HasSuffix(pkg.ImportPath, "github.com/tsaikd/KDGoLib"))
	}

	pkglist, err = ParsePackagePaths("..", "github.com/tsaikd/KDGoLib")
	require.NoError(err)
	require.EqualValues(1, pkglist.Len())
	for pkg := range pkglist.Map() {
		require.True(strings.HasSuffix(pkg.ImportPath, "github.com/tsaikd/KDGoLib"))
	}

	pkglist, err = ParsePackagePaths("..", "github.com/tsaikd/KDGoLib/...")
	require.NoError(err)
	require.True(pkglist.Len() > 1)

	pkglist, err = ParsePackagePaths("..", "./...")
	require.NoError(err)
	require.True(pkglist.Len() > 1)

	pkglist, err = ParsePackagePaths("..", "./cliutil/...")
	require.NoError(err)
	require.True(pkglist.Len() > 1)
}
