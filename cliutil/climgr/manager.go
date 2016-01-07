package climgr

import "github.com/codegangsta/cli"

// Manager manage all registed cli module
type Manager struct {
	moduleNameMap  map[string]Module
	commandNameMap map[string]cli.Command
	flagNameMap    map[string]cli.Flag
}

// NewManager create a new manager instance
func NewManager() *Manager {
	return &Manager{
		moduleNameMap:  map[string]Module{},
		commandNameMap: map[string]cli.Command{},
		flagNameMap:    map[string]cli.Flag{},
	}
}

// Reset clean all registed modules, commands, flags
func (t *Manager) Reset() {
	t.moduleNameMap = map[string]Module{}
	t.commandNameMap = map[string]cli.Command{}
	t.flagNameMap = map[string]cli.Flag{}
}
