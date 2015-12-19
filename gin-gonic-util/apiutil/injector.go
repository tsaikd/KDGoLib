package apiutil

import (
	"net/http"

	"github.com/codegangsta/inject"
	"github.com/gin-gonic/gin"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/injectutil"
)

// errors
var (
	ErrorInjectorNotFound = errutil.ErrorFactory("injector not found, remember to use Injector middleware ?")
)

var (
	metaInjectorKey = "injector"
)

// Handler is a injectable handler
type Handler interface{}

// Wrap handler for injection
func Wrap(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		inj, exist := Get(c)
		if !exist {
			c.AbortWithError(http.StatusServiceUnavailable, ErrorInjectorNotFound.New(nil))
			return
		}

		if _, err := injectutil.Invoke(inj, handler); err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
	}
}

// Injector is gin middleware to support handler injection, use Wrap() for custom handler
func Injector(parent inject.Injector) gin.HandlerFunc {
	return func(c *gin.Context) {
		inj := inject.New()
		inj.SetParent(parent)
		inj.Map(inj)
		inj.Map(c)

		c.Set(metaInjectorKey, inj)

		c.Next()
	}
}

// Get injector from gin context
func Get(c *gin.Context) (inj inject.Injector, exist bool) {
	injvalue, exist := c.Get(metaInjectorKey)
	if exist {
		inj = injvalue.(inject.Injector)
		return injvalue.(inject.Injector), true
	}
	return nil, false
}
