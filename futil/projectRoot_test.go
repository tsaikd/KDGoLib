package futil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchProjectRoot(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	if dir, err := SearchProjectRoot("", []string{"futil"}, []string{".git"}, []string{"KDGoLib"}); assert.NoError(err) {
		require.NotEmpty(dir)
	}
}
