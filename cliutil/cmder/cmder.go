package cmder

import (
	"os"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/version"
	"gopkg.in/urfave/cli.v2"
)

// Main entry point
func Main(module Module, commandModules ...Module) {
	app := NewApp(module)
	app.Version = version.String()

	usedModules := []Module{module}
	for _, mod := range module.Depend.DependOnce() {
		usedModules = append(usedModules, *mod)
	}

	cmds := []*cli.Command{}
	for _, cmdmod := range commandModules {
		cmds = append(cmds, NewCommand(cmdmod, usedModules...))
	}
	if VersionModule != nil {
		cmds = append(cmds, NewCommand(*VersionModule, usedModules...))
	}
	app.Commands = cmds

	errutil.Trace(app.Run(os.Args))
}
