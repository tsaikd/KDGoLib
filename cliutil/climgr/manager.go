package climgr

import (
	"sort"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorModuleRegisted1 = errutil.ErrorFactory("Module already registed: %q")
	ErrorModuleNotFound1 = errutil.ErrorFactory("Module not found: %q")
)

// Manager manage all registed cli module
type Manager struct {
	moduleNameMap map[string]Definition
}

// NewManager create a new manager instance
func NewManager() *Manager {
	return &Manager{
		moduleNameMap: map[string]Definition{},
	}
}

// Add module to manager
func (t *Manager) Add(module Definition) (err error) {
	key := module.Name
	if _, exist := t.moduleNameMap[key]; exist {
		return ErrorModuleRegisted1.New(nil, key)
	}
	t.moduleNameMap[key] = module
	return
}

// Delete module from manager
func (t *Manager) Delete(module Definition) {
	key := module.Name
	delete(t.moduleNameMap, key)
}

// Reset clean all registed modules
func (t *Manager) Reset() {
	t.moduleNameMap = map[string]Definition{}
}

func (t sortByPriorityName) Less(i, j int) bool {
	x := t.Definitions[i]
	y := t.Definitions[j]
	if x.Priority < y.Priority {
		return true
	}
	if x.Priority > y.Priority {
		return false
	}
	return x.Name < y.Name
}

// GetModules return registed modules, sorted by priority, module name
func (t *Manager) GetModules() (results Definitions) {
	for _, module := range t.moduleNameMap {
		results = append(results, module)
	}
	sort.Sort(sortByPriorityName{results})
	return
}

// SelectModules return selected modules, sorted by priority, module name
func (t *Manager) SelectModules(moduleNames ...string) (results Definitions, err error) {
	for _, moduleName := range moduleNames {
		module, exist := t.moduleNameMap[moduleName]
		if !exist {
			return Definitions{}, ErrorModuleNotFound1.New(nil, moduleName)
		}
		results = append(results, module)
	}
	sort.Sort(sortByPriorityName{results})
	return
}
