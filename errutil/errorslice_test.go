package errutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	assert := assert.New(t)

	err1 := New("err1")
	assert.Error(err1)
	assert.Equal("err1", err1.Error())

	err2 := New("err2", err1)
	assert.Error(err2)
	assert.Equal("err2\nerr1", err2.Error())

	err3 := New("err3", err2)
	assert.Error(err3)
	assert.Equal("err3\nerr2\nerr1", err3.Error())

	err4 := New("err4", err2, err1)
	assert.Error(err4)
	assert.Equal("err4\nerr2\nerr1\nerr1", err4.Error())

	data, err := json.Marshal(err3)
	assert.NoError(err)
	assert.Equal(`{"error":"err3","errors":["err3","err2","err1"]}`, string(data))

	unerr := ErrorSlice{}
	err = json.Unmarshal(data, &unerr)
	assert.NoError(err)

	undata, err := json.Marshal(unerr)
	assert.NoError(err)
	assert.Equal(string(undata), string(data))
}
