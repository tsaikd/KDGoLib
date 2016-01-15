package errorJson

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/KDGoLib/errutil"
)

func Test_String(t *testing.T) {
	assert := assert.New(t)

	reserr := &responseError{
		Status:     404,
		ErrorSlice: errutil.NewErrorSlice(errors.New("test error"), errors.New("test error 2")),
	}
	data, err := json.Marshal(reserr)
	assert.NoError(err)
	assert.Contains(string(data), "test error")

	unerr := responseError{}
	err = json.Unmarshal(data, &unerr)
	assert.NoError(err)
	assert.Equal(404, unerr.Status)

	undata, err := json.Marshal(unerr)
	assert.NoError(err)
	assert.Equal(string(undata), string(data))
}
