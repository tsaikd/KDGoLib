package cmder

import "github.com/urfave/cli"

// NewModule create Module instance
func NewModule(name string) *Module {
	return &Module{
		Name:   name,
		Flags:  []cli.Flag{},
		Depend: Modules{},
	}
}

// Module for cli, used for cli.App or cli.Command
type Module struct {
	Name   string
	Usage  string
	Flags  []cli.Flag
	Before cli.ActionFunc
	Action cli.ActionFunc
	After  cli.ActionFunc
	Depend Modules
}

// SetUsage set Module usage
func (t *Module) SetUsage(usage string) *Module {
	t.Usage = usage
	return t
}

// AddFlag add cli.Flag to Module
func (t *Module) AddFlag(flags ...cli.Flag) *Module {
	t.Flags = append(t.Flags, flags...)
	return t
}

// SetBefore set Module before action
func (t *Module) SetBefore(action cli.ActionFunc) *Module {
	t.Before = action
	return t
}

// SetAction set Module action
func (t *Module) SetAction(action cli.ActionFunc) *Module {
	t.Action = action
	return t
}

// SetAfter set Module after action
func (t *Module) SetAfter(action cli.ActionFunc) *Module {
	t.After = action
	return t
}

// AddDepend add dependent module
func (t *Module) AddDepend(modules ...*Module) *Module {
	t.Depend = append(t.Depend, modules...)
	return t
}

// Modules slice of Module
type Modules []*Module

// Flags return all flags in Modules
func (t Modules) Flags(usedModules ...Module) []cli.Flag {
	flags := []cli.Flag{}
	modules := t.DependOnce(usedModules...)
	for _, module := range modules {
		flags = append(flags, module.Flags...)
	}
	return flags
}

// BeforeActions return all before actions in Modules
func (t Modules) BeforeActions(usedModules ...Module) Actions {
	actions := Actions{}
	modules := t.DependOnce(usedModules...)
	for _, module := range modules {
		actions.Add(module.Before)
	}
	return actions
}

// Actions return all actions in Modules
func (t Modules) Actions(usedModules ...Module) Actions {
	actions := Actions{}
	modules := t.DependOnce(usedModules...)
	for _, module := range modules {
		actions.Add(module.Action)
	}
	return actions
}

// AfterActions return all after actions in Modules
func (t Modules) AfterActions(usedModules ...Module) Actions {
	actions := Actions{}
	modules := t.DependOnce(usedModules...)
	for i := len(modules) - 1; i >= 0; i-- {
		module := modules[i]
		actions.Add(module.After)
	}
	return actions
}

// DependOnce return all dependent modules but only show once if depend multiple times
func (t Modules) DependOnce(usedModules ...Module) Modules {
	modules := Modules{}
	usage := nameUsage{}

	for _, module := range usedModules {
		usage.Use(module.Name)
	}

	for _, module := range t {
		depmods := module.Depend.DependOnce()
		for _, depend := range depmods {
			if !usage.IsUsed(depend.Name) {
				usage.Use(depend.Name)
				modules = append(modules, depend)
			}
		}

		if !usage.IsUsed(module.Name) {
			usage.Use(module.Name)
			modules = append(modules, module)
		}
	}

	return modules
}
