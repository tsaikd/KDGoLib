package bqutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	gcKeyFile := ""
	projectID := ""
	invalidOption := ctx
	client, err := NewClient(ctx, gcKeyFile, projectID, invalidOption)
	require.Error(err)
	require.True(ErrUnknownOption2.Match(err))
	require.Nil(client)
}
