package bindstruct

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/martini/errorJson"
)

type apiReq struct {
	ParamA     string `json:"a" valid:"required"`
	ParamB     int64  `json:"b"`
	ParamC     string `json:"c"`
	ParamD     bool   `json:"d"`
	ParamE     bool   `json:"e"`
	ChildSlice []struct {
		Str string `json:"str"`
	} `json:"childslice"`
}

type apiReqRequiredStruct struct {
	RequiredStruct multipart.FileHeader `json:"file" valid:"required"`
}

func Test_binding(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	m := martini.Classic()
	errorJson.BindMartini(m.Martini)

	func() {
		runapi := false
		m.Any("/testsimple", func() {
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testsimple", nil)
		require.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		require.Equal(200, httpRecorder.Code)
		require.Empty(httpRecorder.Body.String())
		require.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testquery", BindStruct(apiReq{}), func(areq apiReq) {
			require.Equal("parama", areq.ParamA)
			require.Equal(int64(123), areq.ParamB)
			require.True(areq.ParamD)
			require.False(areq.ParamE)
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testquery?a=parama&b=123&d=true", nil)
		require.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		require.Equal(200, httpRecorder.Code)
		require.Empty(httpRecorder.Body.String())
		require.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testparam/:a/:b/:d", BindStruct(apiReq{}), func(areq apiReq) {
			require.Equal("parama", areq.ParamA)
			require.Equal(int64(123), areq.ParamB)
			require.True(areq.ParamD)
			require.False(areq.ParamE)
			runapi = true
		})
		req, err := http.NewRequest("POST", "/testparam/parama/123/true", nil)
		require.NoError(err)
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		require.Equal(200, httpRecorder.Code)
		require.Empty(httpRecorder.Body.String())
		require.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/testjsonbody", BindStruct(apiReq{}), func(areq apiReq) {
			require.Equal("parama", areq.ParamA)
			require.Equal(int64(123), areq.ParamB)
			require.Equal("paramc", areq.ParamC)
			require.True(areq.ParamD)
			require.False(areq.ParamE)
			require.Len(areq.ChildSlice, 1)
			require.Equal(areq.ChildSlice[0].Str, "text")
			runapi = true
		})
		body, err := json.Marshal(map[string]interface{}{
			"a": " parama ",
			"b": 123,
			"d": true,
			"childslice": []map[string]interface{}{
				map[string]interface{}{
					"str": "text",
				},
			},
		})
		require.NoError(err)
		req, err := http.NewRequest("POST", "/testjsonbody?c=paramc", bytes.NewReader(body))
		require.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		require.Equal(200, httpRecorder.Code)
		require.Empty(httpRecorder.Body.String())
		require.True(runapi)
	}()

	func() {
		runapi := false
		m.Any("/teststructrequired", BindStruct(apiReqRequiredStruct{}), func(areq apiReqRequiredStruct) {
			runapi = true
		})
		req, err := http.NewRequest("POST", "/teststructrequired", nil)
		require.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		httpRecorder := httptest.NewRecorder()
		m.ServeHTTP(httpRecorder, req)
		require.Equal(404, httpRecorder.Code)
		require.Contains(httpRecorder.Body.String(), "bindStruct.go:97")
		require.False(runapi)
	}()
}
