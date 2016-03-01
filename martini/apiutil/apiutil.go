package apiutil

import (
	"sort"

	"github.com/go-martini/martini"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorAPIRegisted2 = errutil.NewFactory("API already registed for method=%q pattern=%1")
)

var (
	registedApiMethodPatternMap = map[string]apiConfig{}
	registedApiCallerMap        = map[string]apiConfig{}
)

// BindRouter bind all registed apis to martini router
func BindRouter(router martini.Router) {
	for _, apiconf := range registedApiMethodPatternMap {
		router.AddRoute(apiconf.Method, apiconf.Pattern, apiconf.Handlers...)
	}
}

// RegistedApiMap return all registed API
func RegistedApiMap() map[string]apiConfig {
	return registedApiMethodPatternMap
}

// RegistedSortedApiMap return all registed API sorted
func RegistedSortedApiMap() []apiConfig {
	results := []apiConfig{}

	keys := []string{}
	for key := range registedApiCallerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		apiconf := registedApiCallerMap[key]
		results = append(results, apiconf)
	}

	return results
}

// Regist API
func Regist(method string, pattern string, handlers ...martini.Handler) {
	WithInfo(nil).Regist(method, pattern, handlers...)
}

// Reset all registed API, used for testing
func Reset() {
	registedApiMethodPatternMap = map[string]apiConfig{}
	registedApiCallerMap = map[string]apiConfig{}
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
