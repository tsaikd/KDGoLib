package enumutil

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert := assert.New(t)
	assert.NotNil(assert)

	assert.EqualValues(2, TestEnumB)

	enumstr := fmt.Sprintf("%s", TestEnumB)
	assert.Equal("b", enumstr)

	enumjson, err := json.Marshal(TestEnumB)
	assert.NoError(err)
	assert.Equal(`"b"`, string(enumjson))

	enumdata := TestEnumA
	err = json.Unmarshal([]byte(`"b"`), &enumdata)
	assert.NoError(err)
	assert.Equal(TestEnumB, enumdata)
}
