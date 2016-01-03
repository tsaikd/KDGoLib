package apimgr

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorAPIRegisted2 = errutil.ErrorFactory("API already registed for method=%q pattern=%q")
)

// Manager manage all registed APIs
type Manager struct {
	BasePackage         string
	GetMethodPatternKey func(*Manager, Definition) string
	NameGenerator       func(*Manager, Definition) (name string, fullname string)
	MethodGenerator     func(*Manager, Definition) string
	PatternGenerator    func(*Manager, Definition) string
	apiMethodPatternMap map[string]Definition
}

// NewManager create a new manager instance
func NewManager(basePackage interface{}) *Manager {
	return &Manager{
		BasePackage:         getPackagePath(reflect.ValueOf(basePackage)),
		GetMethodPatternKey: getMethodPatternKey,
		NameGenerator:       nameGenerator,
		MethodGenerator:     methodGenerator,
		PatternGenerator:    patternGenerator,
		apiMethodPatternMap: map[string]Definition{},
	}
}

// Add api to manager
func (t *Manager) Add(api Definition) (err error) {
	if api.Name == "" || api.FullName == "" {
		name, fullname := t.NameGenerator(t, api)
		if api.Name == "" {
			api.Name = name
		}
		if api.FullName == "" {
			api.FullName = fullname
		}
	}
	if api.Method == "" {
		api.Method = t.MethodGenerator(t, api)
	}
	if api.Pattern == "" {
		api.Pattern = t.PatternGenerator(t, api)
	}
	key := t.GetMethodPatternKey(t, api)
	if _, exist := t.apiMethodPatternMap[key]; exist {
		return ErrorAPIRegisted2.New(nil, api.Method, api.Pattern)
	}
	t.apiMethodPatternMap[key] = api
	return
}

// Delete api from manager
func (t *Manager) Delete(api Definition) {
	key := t.GetMethodPatternKey(t, api)
	delete(t.apiMethodPatternMap, key)
}

// Reset clean all registed apis
func (t *Manager) Reset() {
	t.apiMethodPatternMap = map[string]Definition{}
}

// GetMethodPatternMap return string definition map, key := method pattern
func (t *Manager) GetMethodPatternMap() map[string]Definition {
	return t.apiMethodPatternMap
}

// GetSortedAPIsByPkgPath return all registed APIs sorted by package path
// NOTE: this is slow for sorting in runtime
func (t *Manager) GetSortedAPIsByPkgPath() (results []Definition) {
	pkgpathMap := map[string]Definition{}
	keys := []string{}

	for _, api := range t.apiMethodPatternMap {
		pkgpath := getPackagePath(reflect.ValueOf(api.Request))
		key := pkgpath + " " + t.GetMethodPatternKey(t, api)
		pkgpathMap[key] = api
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		results = append(results, pkgpathMap[key])
	}

	return results
}

func getMethodPatternKey(manager *Manager, api Definition) string {
	return api.Method + " " + api.Pattern
}

func nameGenerator(manager *Manager, api Definition) (name string, fullname string) {
	pkgpath := getPackagePath(reflect.ValueOf(api.Request))
	relpath, _ := filepath.Rel(manager.BasePackage, pkgpath)
	relpaths := strings.Split(relpath, string(os.PathSeparator))
	lenRelpaths := len(relpaths)
	if lenRelpaths > 0 {
		name += camelcase(relpaths[lenRelpaths-1])
		for i := 0; i < lenRelpaths-1; i++ {
			name += camelcase(relpaths[i])
		}
	}
	fullname = name
	if api.Version > 0 {
		fullname += "V" + strconv.FormatUint(uint64(api.Version), 10)
	}
	return
}

func methodGenerator(manager *Manager, api Definition) string {
	return "GET"
}

func patternGenerator(manager *Manager, api Definition) string {
	pattern := "/"
	if api.Version > 0 {
		pattern += strconv.FormatUint(uint64(api.Version), 10) + "/"
	}
	pkgpath := getPackagePath(reflect.ValueOf(api.Request))
	relpath, _ := filepath.Rel(manager.BasePackage, pkgpath)
	relpaths := strings.Split(relpath, string(os.PathSeparator))
	lenRelpaths := len(relpaths)
	if lenRelpaths > 0 {
		pattern += relpaths[lenRelpaths-1]
		for i := 0; i < lenRelpaths-1; i++ {
			pattern += "/" + relpaths[i]
		}
	}
	return pattern
}

func camelcase(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}
