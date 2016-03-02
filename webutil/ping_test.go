package webutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ping(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := PingIgnoreCertificate("invalid url")
	require.Error(err)
}
