package httpRender_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/webutil/render"
	"github.com/tsaikd/KDGoLib/webutil/render/httpRender"
)

func TestRender(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	req := httptest.NewRequest("GET", "/", nil)
	lastModified := time.Now()
	req.Header.Set(httpRender.HeaderIfModifiedSince, lastModified.Format(http.TimeFormat))

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.JSON(nil)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("application/json; charset=UTF-8", w.Header().Get(httpRender.HeaderContentType))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.JSON(nil)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("null\n", string(r.GetBody()))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.Error(errutil.New("error_message"))
		require.NotNil(r.GetError())
		require.Equal(http.StatusNotFound, r.GetStatus())
		require.Contains(string(r.GetBody()), `"error":"error_message"`)
		require.Contains(string(r.GetBody()), `"status":404`)
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.Error(nil)
		require.Nil(r.GetError())
		require.Equal(0, r.GetStatus())
		require.Equal("", string(r.GetBody()))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.Redirect(http.StatusSeeOther, "http://example.com")
		require.Nil(r.GetError())
		require.Equal(http.StatusSeeOther, r.GetStatus())
		require.Equal("", string(r.GetBody()))
		require.Equal("http://example.com", w.Header().Get("Location"))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.SetCookie(&http.Cookie{
			Name:     "testcookie",
			Value:    "test cookie value",
			Domain:   "localhost",
			HttpOnly: true,
		})
		if cookie := r.GetResponseHeader().Get("Set-Cookie"); assert.NotEmpty(cookie) {
			// in go dev: cookie will be `testcookie="test cookie value"; Domain=localhost; HttpOnly`
			// require.Equal("testcookie=test cookie value; Domain=localhost; HttpOnly", cookie)
			require.Contains(cookie, "testcookie=")
			require.Contains(cookie, "test cookie value")
			require.Contains(cookie, "Domain=localhost; HttpOnly")
		}
		r.JSON(nil)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("null\n", string(r.GetBody()))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		header := http.Header{}
		header.Set("custom-header", "custom header value")
		r.WriteResponse(header, http.StatusOK, []byte("test bytes body"))
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal([]byte("test bytes body"), r.GetBody())
		require.Equal("custom header value", w.Header().Get("custom-header"))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		writer := r.GetResponseWriter()
		writer.Header().Set("custom-header", "custom header value")
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("part1"))
		require.NoError(err)
		_, err = writer.Write([]byte("part2"))
		require.NoError(err)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("custom header value", w.Header().Get("custom-header"))
		require.Equal([]byte("part1part2"), r.GetBody())
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		require.True(r.LastModified(lastModified))
		require.Nil(r.GetError())
		require.Equal(http.StatusNotModified, r.GetStatus())
		require.EqualValues(0, r.GetSize())
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.SetCacheControl(render.CacheControlNoStore, 0)
		r.JSON(nil)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("no-store", w.Header().Get(httpRender.HeaderCacheControl))
	}

	if w := httptest.NewRecorder(); assert.NotNil(w) {
		r := httpRender.New(w, req)
		r.SetCacheControl(render.CacheControlPublic, 600*time.Second)
		r.JSON(nil)
		require.Nil(r.GetError())
		require.Equal(http.StatusOK, r.GetStatus())
		require.Equal("public, max-age=600", w.Header().Get(httpRender.HeaderCacheControl))
	}
}
