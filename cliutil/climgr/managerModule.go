package climgr

import (
	"sort"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorModuleRegisted1 = errutil.NewFactory("Module already registed: %q")
	ErrorModuleNotFound1 = errutil.NewFactory("Module not found: %q")
)

// AddModule add module to manager
func (t *Manager) AddModule(module Module) (err error) {
	key := module.Name
	if _, exist := t.moduleNameMap[key]; exist {
		return ErrorModuleRegisted1.New(nil, key)
	}
	t.moduleNameMap[key] = module
	return
}

// DeleteModule delete module from manager
func (t *Manager) DeleteModule(module Module) {
	key := module.Name
	delete(t.moduleNameMap, key)
}

// GetModules return registed modules, sorted by priority, module name
func (t *Manager) GetModules() (results Modules) {
	for _, module := range t.moduleNameMap {
		results = append(results, module)
	}
	sort.Sort(sortModulesByPriorityName{results})
	return
}

// SelectModules return selected modules, sorted by priority, module name
func (t *Manager) SelectModules(moduleNames ...string) (results Modules, err error) {
	for _, moduleName := range moduleNames {
		module, exist := t.moduleNameMap[moduleName]
		if !exist {
			return Modules{}, ErrorModuleNotFound1.New(nil, moduleName)
		}
		results = append(results, module)
	}
	sort.Sort(sortModulesByPriorityName{results})
	return
}
