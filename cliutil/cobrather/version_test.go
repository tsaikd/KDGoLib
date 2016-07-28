package cobrather_test

import (
	"testing"

	semver "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func Test_version(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	ver, err := semver.NewVersion("1.2.3")
	require.NoError(err)

	testdata := []struct {
		Valid bool
		Range string
	}{
		{true, ">0"},
		{true, ">=0"},
		{true, ">=0.0"},
		{true, ">=0.0.0"},
		{true, ">=0.0.1"},
		{true, ">=1.2.3"},
		{true, ">=1.2.3-alpha"},
		{true, ">=1.2.3-beta"},
		{true, ">=1.2.3-dev"},
		{true, ">=1.2 , <2"},
		{true, ">1.2"},
		{true, ">1.2.2"},
		{false, ">1.2.3"},
		{false, ">=1.3"},
		{false, ">=1.3 , <2"},
		{false, ">=2"},
	}

	for _, data := range testdata {
		verRange, err := semver.NewConstraint(data.Range)
		require.NoError(err)
		require.Equal(data.Valid, verRange.Check(ver), data.Range)
	}
}
