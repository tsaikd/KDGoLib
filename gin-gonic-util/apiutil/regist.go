package apiutil

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorAPIRegisted2 = errutil.ErrorFactory("API already registed for method=%q path=%1")
)

var (
	registedAPIMethodPathMap = map[string]apiConfig{}
	registedAPICallerMap     = map[string]apiConfig{}
)

type apiConfig struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
	Info     interface{}
}

type apiOrm struct {
	Info interface{}
}

// WithInfo return API ORM pointer for regist API
func WithInfo(info interface{}) *apiOrm {
	return &apiOrm{
		Info: info,
	}
}

func getMethodPath(method string, path string) string {
	return method + " " + path
}

func getCallerKey() string {
	var (
		thisfile   string
		callerfile string
		line       int
		ok         bool
	)
	if _, thisfile, _, ok = runtime.Caller(0); ok {
		thisdir := path.Dir(thisfile)
		for i := 1; ok; i++ {
			_, callerfile, line, ok = runtime.Caller(i)
			if !strings.HasPrefix(callerfile, thisdir) {
				return callerfile + ":" + strconv.Itoa(line)
			}
		}
	}
	return ""
}

func (t *apiOrm) Regist(method string, path string, handlers ...gin.HandlerFunc) {
	methodPattern := getMethodPath(method, path)
	callerKey := getCallerKey()
	if callerKey == "" {
		callerKey = methodPattern
	}

	_, ok := registedAPIMethodPathMap[methodPattern]
	if ok {
		panic(ErrorAPIRegisted2.New(nil, method, path))
	}

	apiconf := apiConfig{
		Method:   method,
		Path:     path,
		Handlers: handlers,
		Info:     t.Info,
	}

	registedAPIMethodPathMap[methodPattern] = apiconf
	registedAPICallerMap[callerKey] = apiconf

	return
}

func (t *apiOrm) POST(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodPOST.String(), path, handlers...)
}

func (t *apiOrm) GET(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodGET.String(), path, handlers...)
}

func (t *apiOrm) DELETE(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodDELETE.String(), path, handlers...)
}

func (t *apiOrm) PATCH(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodPATCH.String(), path, handlers...)
}

func (t *apiOrm) PUT(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodPUT.String(), path, handlers...)
}

func (t *apiOrm) OPTIONS(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodOPTIONS.String(), path, handlers...)
}

func (t *apiOrm) HEAD(path string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodHEAD.String(), path, handlers...)
}

func (t *apiOrm) Any(pattern string, handlers ...gin.HandlerFunc) {
	t.Regist(MethodAny.String(), pattern, handlers...)
}
