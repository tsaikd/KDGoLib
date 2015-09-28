package bindstruct

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/KDGoLib/martini/errorJson"
)

type apiReq struct {
	ParamA string `json:"a" valid:"required"`
	ParamB int64  `json:"b"`
	ParamC string `json:"c"`
	ParamD bool   `json:"d"`
	ParamE bool   `json:"e"`
}

type apiReqRequiredStruct struct {
	RequiredStruct multipart.FileHeader `json:"file" valid:"required"`
}

func Test_binding(t *testing.T) {
	assert := assert.New(t)

	m := martini.Classic()
	m.Map(errorJson.ReturnErrorProvider())
	m.Use(render.Renderer())

	func() {
		runapi := false
		m.Any("/testsimple", func() {
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testsimple", nil)
		assert.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		assert.Empty(httpRecorder.Body.String())
		assert.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testquery", BindStruct(apiReq{}), func(areq apiReq) {
			assert.Equal("parama", areq.ParamA)
			assert.Equal(123, areq.ParamB)
			assert.True(areq.ParamD)
			assert.False(areq.ParamE)
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testquery?a=parama&b=123&d=true", nil)
		assert.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		assert.Empty(httpRecorder.Body.String())
		assert.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testparam/:a/:b/:d", BindStruct(apiReq{}), func(areq apiReq) {
			assert.Equal("parama", areq.ParamA)
			assert.Equal(123, areq.ParamB)
			assert.True(areq.ParamD)
			assert.False(areq.ParamE)
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testparam/parama/123/true", nil)
		assert.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		assert.Empty(httpRecorder.Body.String())
		assert.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testjsonbody", BindStruct(apiReq{}), func(areq apiReq) {
			assert.Equal("parama", areq.ParamA)
			assert.Equal(123, areq.ParamB)
			assert.Equal("paramc", areq.ParamC)
			assert.True(areq.ParamD)
			assert.False(areq.ParamE)
			runapi = true
		})
		body, err := json.Marshal(map[string]interface{}{
			"a": " parama ",
			"b": 123,
			"d": true,
		})
		assert.NoError(err)
		req, err := http.NewRequest("POST", "/testjsonbody?c=paramc", bytes.NewReader(body))
		assert.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		assert.Empty(httpRecorder.Body.String())
		assert.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/teststructrequired", BindStruct(apiReqRequiredStruct{}), func(areq apiReqRequiredStruct) {
			runapi = true
		})
		req, err := http.NewRequest("POST", "/teststructrequired", nil)
		assert.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		assert.NotEmpty(httpRecorder.Body.String())
		assert.False(runapi)
	}()
}
