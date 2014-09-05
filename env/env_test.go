package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetString(t *testing.T) {
	assert := assert.New(t)

	value := GetString("IMPOSSIBLE_ENV_KEY", "Impossible value !")
	assert.Equal("Impossible value !", value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	path := os.Getenv("PATH")
	value = GetString("PATH", "")
	assert.Equal(path, value, "Get ENV 'PATH'")
}
