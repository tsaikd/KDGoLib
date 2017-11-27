package gofmtutil_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/gofmtutil"
)

func TestGoImports(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	if fmted, err := gofmtutil.GoImports([]byte("")); assert.NoError(err) {
		require.Len(fmted, 0)
	}

	if fmted, err := gofmtutil.GoImports([]byte(`
package test

func main () {
fmt.Println( "test", context.Background ( ) )
}

	`)); assert.NoError(err) {
		require.Equal(strings.TrimSpace(`
package test

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("test", context.Background())
}
		`)+"\n", string(fmted))
	}
}
