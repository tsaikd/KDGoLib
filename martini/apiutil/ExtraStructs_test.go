package apiutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtraStructs(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	extraStructs := ExtraStructs{
		&ExtraStruct{
			FieldName: "a",
		},
	}
	assert.Len(extraStructs, 1)

	extraStructs.Upsert("a", reflect.TypeOf(0))
	assert.Len(extraStructs, 1)

	extraStructs.Upsert("b", reflect.TypeOf(0))
	assert.Len(extraStructs, 2)
}
