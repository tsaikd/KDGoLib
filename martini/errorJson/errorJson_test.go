package errorJson

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	assert := assert.New(t)

	reserr := newResponseError(404, errors.New("test error"), errors.New("test error 2"))
	data, err := json.Marshal(reserr)
	assert.NoError(err)
	assert.Contains(string(data), "test error")

	unerr := ResponseError{}
	err = json.Unmarshal(data, &unerr)
	assert.NoError(err)
	assert.Equal(404, unerr.Status)

	undata, err := json.Marshal(unerr)
	assert.NoError(err)
	assert.Equal(string(undata), string(data))
}
