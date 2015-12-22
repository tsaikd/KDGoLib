package apiutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RequestParams(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	reqparams := RequestParams{
		RequestParam{
			FieldName: "a",
		},
	}
	assert.Len(reqparams, 1)

	err := reqparams.Delete("a")
	assert.NoError(err)
	assert.Len(reqparams, 0)

	reqparams.Upsert("b", nil)
	assert.Len(reqparams, 1)
	assert.Equal("", reqparams[0].FieldType.Name)

	reqparams.Upsert("b", &reflect.StructField{
		Name: "b",
	})
	assert.Len(reqparams, 1)
	assert.Equal("b", reqparams[0].FieldType.Name)

	cloneparams := reqparams.Clone()
	cloneparams.Delete("b")
	assert.Len(reqparams, 1)
}
