package climgr

import "github.com/codegangsta/cli"

// ModuleFunc module function for execute
type ModuleFunc func(c *cli.Context) (err error)

// Definition contain module detail
type Definition struct {
	Name       string
	ModuleFunc ModuleFunc
	CloseFunc  ModuleFunc
	Optional   bool
	Priority   int8 // smaller value means high priority
}

// Definitions slice of Definitions
type Definitions []Definition

// Len return length of slice
func (t Definitions) Len() int {
	return len(t)
}

// Swap change slice 2 element position
func (t Definitions) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type sortByPriorityName struct {
	Definitions
}
