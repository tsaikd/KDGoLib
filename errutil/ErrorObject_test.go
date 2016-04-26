package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorObject_New(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := New("new error")
	require.Error(err)
	require.Equal("ErrorObject_test.go", err.FileName())
	require.Equal(14, err.Line())
	require.Equal("new error", err.Error())
	require.Nil(err.Factory())
	require.Nil(err.Parent())
	require.Equal(1, Length(err))
}

func Test_ErrorObject_NewErrors(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := NewErrors(errors.New("new errors"))
	require.Error(err)
	require.Equal("ErrorObject_test.go", err.FileName())
	require.Equal(28, err.Line())
	require.Equal("new errors", err.Error())
	require.Nil(err.Factory())
	require.Nil(err.Parent())
	require.Equal(1, Length(err))
}

func Test_ErrorObject_chain_basic(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err1 := New("err1")
	require.Error(err1)
	require.Equal("err1", err1.Error())
	require.Nil(err1.Factory())
	require.Nil(err1.Parent())
	require.Equal(1, Length(err1))

	err2 := New("err2", err1)
	require.Error(err2)
	require.Equal("err2; err1", err2.Error())
	require.Error(err2.Parent())
	require.Equal("err1", err2.Parent().Error())
	require.Equal(2, Length(err2))

	err3 := New("err3", err2)
	require.Error(err3)
	require.Equal("err3; err2; err1", err3.Error())
	require.Error(err3.Parent())
	require.Equal("err2; err1", err3.Parent().Error())
	require.Error(err3.Parent().Parent())
	require.Equal("err1", err3.Parent().Parent().Error())
	require.Equal(3, Length(err3))
}

func Test_ErrorObject_chain_duplicate1(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err1 := New("err1")
	err2 := New("err2", err1)

	err4 := New("err4", err2, err1)
	require.Error(err4)
	require.Equal("err4; err2; err1", err4.Error())
	require.Equal(3, Length(err4))
	data, err := MarshalJSON(err4)
	require.NoError(err)
	require.Contains(string(data), `"error":"err4","errors":["err4","err2","err1"]`)
}

func Test_ErrorObject_chain_duplicate2(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err1 := New("err1")
	err2 := New("err2", err1)
	err3 := New("err3", err2)

	err5 := New("err5", err2, err3)
	require.Error(err5)
	require.Equal("err5; err2; err1; err3", err5.Error())
	require.Equal(-1, Length(err5))
	_, err := MarshalJSON(err5)
	require.Error(err)
	require.True(ErrorWalkLoop.Match(err))
	require.True(ErrorWalkLoop.In(err))
}

func Test_ErrorObject_chain_duplicate3(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err1 := New("err1")
	err2 := New("err2", err1)

	err6 := New("err6", err1, err2)
	require.Error(err6)
	require.Equal("err6; err1; err2", err6.Error())
	require.Equal(-1, Length(err6))
	_, err := MarshalJSON(err6)
	require.Error(err)
	require.True(ErrorWalkLoop.Match(err))
	require.True(ErrorWalkLoop.In(err))
}
