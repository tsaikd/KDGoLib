package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetString(t *testing.T) {
	var value string
	assert := assert.New(t)

	value = GetString("IMPOSSIBLE_ENV_KEY", "Impossible value !")
	assert.Equal("Impossible value !", value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	path := os.Getenv("PATH")
	value = GetString("PATH", "")
	assert.Equal(path, value, "Get ENV 'PATH'")
}

func Test_GetBool(t *testing.T) {
	var value bool
	assert := assert.New(t)

	value = GetBool("IMPOSSIBLE_ENV_KEY", true)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	value = GetBool("IMPOSSIBLE_ENV_KEY", false)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_TRUE", "true")
	value = GetBool("IMPOSSIBLE_ENV_KEY_TRUE", true)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_TRUE'")

	value = GetBool("IMPOSSIBLE_ENV_KEY_TRUE", false)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_TRUE'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_FALSE", "false")
	value = GetBool("IMPOSSIBLE_ENV_KEY_FALSE", true)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_FALSE'")

	value = GetBool("IMPOSSIBLE_ENV_KEY_FALSE", false)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_FALSE'")

}
