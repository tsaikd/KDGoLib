package errorJson

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_RenderErrorJSON(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	defer func() {
		gin.SetMode(mode)
	}()

	router := gin.New()
	router.Use(
		RenderErrorJSON,
	)
	router.GET("/error", func(c *gin.Context) {
		c.Error(errors.New("test error"))
	})
	router.GET("/error2", func(c *gin.Context) {
		c.Error(errors.New("test error 1"))
		c.Error(errors.New("test error 2"))
	})

	w := performRequest(router, "GET", "/error")
	assert.Equal(http.StatusNotFound, w.Code)
	assert.Contains(w.Body.String(), "test error")

	w = performRequest(router, "GET", "/error2")
	assert.Equal(http.StatusNotFound, w.Code)
	assert.Contains(w.Body.String(), "test error 1")
	assert.Contains(w.Body.String(), "test error 2")
}
