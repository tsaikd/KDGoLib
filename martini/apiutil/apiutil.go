package apiutil

import (
	"sort"

	"github.com/go-martini/martini"
	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	ErrorAPIRegisted2 = errutil.ErrorFactory("API already registed for method=%q pattern=%1")
)

var (
	registedApiMethodPatternMap = map[string]apiConfig{}
	registedApiCallerMap        = map[string]apiConfig{}
)

// bind all registed apis to martini router
func BindRouter(router martini.Router) {
	for _, apiconf := range registedApiMethodPatternMap {
		router.AddRoute(apiconf.Method, apiconf.Pattern, apiconf.Handlers...)
	}
}

func RegistedApiMap() map[string]apiConfig {
	return registedApiMethodPatternMap
}

func RegistedSortedApiMap() []apiConfig {
	results := []apiConfig{}

	keys := []string{}
	for key, _ := range registedApiCallerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		apiconf := registedApiCallerMap[key]
		results = append(results, apiconf)
	}

	return results
}

func Regist(method string, pattern string, handlers ...martini.Handler) {
	WithInfo(nil).Regist(method, pattern, handlers...)
}

func Get(pattern string, handlers ...martini.Handler) {
	Regist("GET", pattern, handlers...)
}

func Patch(pattern string, handlers ...martini.Handler) {
	Regist("PATCH", pattern, handlers...)
}

func Post(pattern string, handlers ...martini.Handler) {
	Regist("POST", pattern, handlers...)
}

func Put(pattern string, handlers ...martini.Handler) {
	Regist("PUT", pattern, handlers...)
}

func Delete(pattern string, handlers ...martini.Handler) {
	Regist("DELETE", pattern, handlers...)
}

func Options(pattern string, handlers ...martini.Handler) {
	Regist("OPTIONS", pattern, handlers...)
}

func Head(pattern string, handlers ...martini.Handler) {
	Regist("HEAD", pattern, handlers...)
}

func Any(pattern string, handlers ...martini.Handler) {
	Regist("*", pattern, handlers...)
}
