package requireutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RequireJSONField(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	testdata := []byte(`{"foo":"foo text","bar":"bar text","nilField":null,"number":1}`)
	RequireJSONFieldEqual(t, "foo text", testdata, "foo")
	RequireJSONFieldNotEqual(t, "bar text", testdata, "foo")
	RequireJSONFieldEqualValues(t, int64(1), testdata, "number")
	RequireJSONFieldEqualValues(t, int8(1), testdata, "number")
	RequireJSONFieldExist(t, testdata, "foo")
	RequireJSONFieldNotExist(t, testdata, "notExistField")
	RequireJSONFieldNil(t, testdata, "nilField")
	RequireJSONFieldNotNil(t, testdata, "bar")
}
