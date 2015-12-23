package apiutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AnalyzeRequestStruct(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	reqparams := AnalyzeRequestStruct(
		"/1/login/:username",
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Inherit  struct {
				Name string `json:"name"`
			}
		}{},
		nil,
		nil,
	)
	assert.Equal("LoginV1", reqparams.FuncName)
	assert.Len(reqparams.FuncArgs, 2)
	assert.Len(reqparams.Params, 1)
	assert.Len(reqparams.ExtraStructs, 0)
	assert.Equal("/1/login/username", reqparams.Path)

	reqparams = AnalyzeRequestStruct(
		"/2/login/:username",
		struct {
			Name string `json:"name"`
		}{},
		nil,
		nil,
	)
	assert.Equal("LoginV2", reqparams.FuncName)
	assert.Len(reqparams.FuncArgs, 2)
	assert.Len(reqparams.Params, 1)
	assert.Len(reqparams.ExtraStructs, 0)
	assert.Equal("/2/login/username", reqparams.Path)

	reqparams = AnalyzeRequestStruct(
		"/3/login/:username",
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Inherit  struct {
				Name string `json:"name"`
			} `apijs:"inherit"`
		}{},
		nil,
		nil,
	)
	assert.Equal("LoginV3", reqparams.FuncName)
	assert.Len(reqparams.FuncArgs, 3)
	assert.Len(reqparams.Params, 2)
	assert.Len(reqparams.ExtraStructs, 0)
	assert.Equal("/3/login/username", reqparams.Path)

	reqparams = AnalyzeRequestStruct(
		"/4/login/:username",
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Extra    struct {
				Age     int `json:"age"`
				Address struct {
					Country int `json:"country"`
				} `json:"address"`
				Phones []struct {
					Code   string `json:"code"`
					Number string `json:"number"`
				} `json:"phones"`
			} `json:"extra"`
		}{},
		nil,
		nil,
	)
	assert.Equal("LoginV4", reqparams.FuncName)
	assert.Len(reqparams.FuncArgs, 3)
	assert.Len(reqparams.Params, 2)
	assert.Len(reqparams.ExtraStructs, 3)
	if len(reqparams.ExtraStructs) > 0 {
		extraStruct := reqparams.ExtraStructs[0]
		assert.Equal("extra", extraStruct.FieldName)
		assert.Len(extraStruct.Params, 3)
	}
	if len(reqparams.ExtraStructs) > 1 {
		extraStruct := reqparams.ExtraStructs[1]
		assert.Equal("address", extraStruct.FieldName)
		assert.Len(extraStruct.Params, 1)
	}
	if len(reqparams.ExtraStructs) > 2 {
		extraStruct := reqparams.ExtraStructs[2]
		assert.Equal("phones", extraStruct.FieldName)
		assert.Len(extraStruct.Params, 2)
	}
	assert.Equal("/4/login/username", reqparams.Path)
}

func Test_GetFuncNameByPattern(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	funcName := GetFuncNameByPattern("/1/version", nil)
	assert.Equal("VersionV1", funcName)

	funcName = GetFuncNameByPattern("/2/users/:user/gists", nil)
	assert.Equal("UsersGistsV2", funcName)
}
