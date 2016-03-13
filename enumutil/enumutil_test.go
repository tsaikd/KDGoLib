package enumutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestEnum int8

const (
	TestEnumA TestEnum = 1 + iota
	TestEnumB
	TestEnumC
)

var testEnum = NewEnumFactory().
	Add(TestEnumA, "a").
	Add(TestEnumB, "b").
	Add(TestEnumC, "c").
	Build()

func (t TestEnum) String() string {
	return testEnum.String(t)
}

func (t TestEnum) MarshalJSON() ([]byte, error) {
	return testEnum.MarshalJSON(t)
}

func (t *TestEnum) UnmarshalJSON(b []byte) (err error) {
	return testEnum.UnmarshalJSON(t, b)
}

func Test_enum(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	require.Equal(int64(2), int64(TestEnumB))

	enumstr := fmt.Sprintf("%s", TestEnumB)
	require.Equal("b", enumstr)

	enumjson, err := json.Marshal(TestEnumB)
	require.NoError(err)
	require.Equal(`"b"`, string(enumjson))

	enumdata := TestEnumA
	err = json.Unmarshal([]byte(`"b"`), &enumdata)
	require.NoError(err)
	require.Equal(TestEnumB, enumdata)

	count := 0
	err = testEnum.Each(func(enum interface{}) (stop bool, err error) {
		count++
		return false, nil
	})
	require.NoError(err)
	require.Equal(3, count)

	count = 0
	err = testEnum.Each(func(enum interface{}) (stop bool, err error) {
		count++
		if count > 1 {
			return true, nil
		}
		return false, nil
	})
	require.NoError(err)
	require.Equal(2, count)

	count = 0
	err = testEnum.Each(func(enum interface{}) (stop bool, err error) {
		count++
		if count > 1 {
			return true, errors.New("")
		}
		return false, nil
	})
	require.Error(err)
	require.Equal(2, count)
}
