package climgr

import "github.com/codegangsta/cli"

// ModuleFunc module function for execute
type ModuleFunc func(c *cli.Context) (err error)

// Module contain module detail
type Module struct {
	Name       string
	ModuleFunc ModuleFunc
	CloseFunc  ModuleFunc
	Optional   bool
	Priority   int8 // smaller value means high priority
}

// Modules slice of Module
type Modules []Module

// Len return length of slice
func (t Modules) Len() int {
	return len(t)
}

// Swap change slice 2 element position
func (t Modules) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type sortModulesByPriorityName struct {
	Modules
}

func (t sortModulesByPriorityName) Less(i, j int) bool {
	x := t.Modules[i]
	y := t.Modules[j]
	if x.Priority < y.Priority {
		return true
	}
	if x.Priority > y.Priority {
		return false
	}
	return x.Name < y.Name
}
