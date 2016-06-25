package cmder

import "gopkg.in/urfave/cli.v2"

// NewCommand return cli.Command with module
func NewCommand(module Module, usedModules ...Module) *cli.Command {
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

	return &cli.Command{
		Name:   module.Name,
		Usage:  module.Usage,
		Flags:  flags,
		Before: cli.BeforeFunc(beforeActions.Wrap()),
		Action: actions.Wrap(),
		After:  cli.AfterFunc(afterActions.Wrap()),
	}
}

// NewApp return *cli.App with module
func NewApp(module Module) *cli.App {
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

	app := cli.NewApp()
	app.Name = module.Name
	app.Usage = module.Usage
	app.Flags = flags
	app.Before = cli.BeforeFunc(beforeActions.Wrap())
	app.Action = actions.WrapMain(module.Action)
	app.After = cli.AfterFunc(afterActions.Wrap())
	return app
}
