package cmder

import (
	"github.com/tsaikd/KDGoLib/version"
	"github.com/urfave/cli"
)

// NewCommand return cli.Command with module
func NewCommand(module Module, usedModules ...Module) cli.Command {
	flags := []cli.Flag{}
	flags = append(flags, module.Flags...)
	flags = append(flags, module.Depend.Flags(usedModules...)...)

	beforeActions := Actions{}
	beforeActions.Add(module.Depend.BeforeActions()...)
	beforeActions.Add(module.Before)
	actions := Actions{}
	actions.Add(module.Depend.Actions()...)
	actions.Add(module.Action)
	afterActions := Actions{}
	afterActions.Add(module.After)
	afterActions.Add(module.Depend.AfterActions()...)

	return cli.Command{
		Name:   module.Name,
		Usage:  module.Usage,
		Flags:  flags,
		Before: cli.BeforeFunc(beforeActions.Wrap()),
		Action: actions.Wrap(),
		After:  cli.AfterFunc(afterActions.Wrap()),
	}
}

// NewApp return *cli.App with module and command modules
func NewApp(module Module, commandModules ...Module) *cli.App {
	flags := []cli.Flag{}
	flags = append(flags, module.Flags...)
	flags = append(flags, module.Depend.Flags()...)

	beforeActions := Actions{}
	beforeActions.Add(module.Depend.BeforeActions()...)
	beforeActions.Add(module.Before)
	actions := module.Depend.Actions()
	afterActions := Actions{}
	afterActions.Add(module.After)
	afterActions.Add(module.Depend.AfterActions()...)

	app := &cli.App{
		Name:    module.Name,
		Usage:   module.Usage,
		Flags:   flags,
		Before:  cli.BeforeFunc(beforeActions.Wrap()),
		Action:  actions.WrapMain(module.Action),
		After:   cli.AfterFunc(afterActions.Wrap()),
		Version: version.String(),
	}

	usedModules := []Module{module}
	for _, mod := range module.Depend.DependOnce() {
		usedModules = append(usedModules, *mod)
	}

	cmds := []cli.Command{}
	for _, cmdmod := range commandModules {
		cmds = append(cmds, NewCommand(cmdmod, usedModules...))
	}
	if VersionModule != nil {
		cmds = append(cmds, NewCommand(*VersionModule, usedModules...))
	}
	app.Commands = cmds

	return app
}
