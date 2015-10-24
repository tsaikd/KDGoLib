package apiutil

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
)

type apiConfig struct {
	Method   string
	Pattern  string
	Handlers []martini.Handler
	Info     interface{}
}

type apiOrm struct {
	Info interface{}
}

func WithInfo(info interface{}) *apiOrm {
	return &apiOrm{
		Info: info,
	}
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

func (t *apiOrm) Regist(method string, pattern string, handlers ...martini.Handler) {
	methodPattern := method + " " + pattern
	callerKey := getCallerKey()
	if callerKey == "" {
		callerKey = methodPattern
	}

	_, ok := registedApiMethodPatternMap[methodPattern]
	if ok {
		panic(ErrorAPIRegisted2.New(nil, method, pattern))
	}

	apiconf := apiConfig{
		Method:   method,
		Pattern:  pattern,
		Handlers: handlers,
		Info:     t.Info,
	}

	registedApiMethodPatternMap[methodPattern] = apiconf
	registedApiCallerMap[callerKey] = apiconf

	return
}

func (t *apiOrm) Get(pattern string, handlers ...martini.Handler) {
	t.Regist("GET", pattern, handlers...)
}

func (t *apiOrm) Patch(pattern string, handlers ...martini.Handler) {
	t.Regist("PATCH", pattern, handlers...)
}

func (t *apiOrm) Post(pattern string, handlers ...martini.Handler) {
	t.Regist("POST", pattern, handlers...)
}

func (t *apiOrm) Put(pattern string, handlers ...martini.Handler) {
	t.Regist("PUT", pattern, handlers...)
}

func (t *apiOrm) Delete(pattern string, handlers ...martini.Handler) {
	t.Regist("DELETE", pattern, handlers...)
}

func (t *apiOrm) Options(pattern string, handlers ...martini.Handler) {
	t.Regist("OPTIONS", pattern, handlers...)
}

func (t *apiOrm) Head(pattern string, handlers ...martini.Handler) {
	t.Regist("HEAD", pattern, handlers...)
}

func (t *apiOrm) Any(pattern string, handlers ...martini.Handler) {
	t.Regist("*", pattern, handlers...)
}
