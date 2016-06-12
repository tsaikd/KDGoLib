package tcmder

import (
	"flag"

	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"gopkg.in/urfave/cli.v2"
)

// NewTest create TestModule instance
func NewTest(module cmder.Module) TestModule {
	return TestModule{
		command: cmder.NewCommand(module),
	}
}

// TestModule used for testing cmder.Module
type TestModule struct {
	command *cli.Command
	context *cli.Context
}

// Setup run before, action in module
func (t TestModule) Setup() (err error) {
	set := flag.NewFlagSet(t.command.Name, flag.PanicOnError)
	for _, f := range t.command.Flags {
		f.Apply(set)
	}

	t.context = cli.NewContext(nil, set, nil)

	if t.command.Before != nil {
		if err = t.command.Before(t.context); err != nil {
			return
		}
	}
	if t.command.Action != nil {
		if err = t.command.Action(t.context); err != nil {
			return
		}
	}

	return
}

// Teardown run after in module
func (t TestModule) Teardown() (err error) {
	if t.command.After != nil {
		if err = t.command.After(t.context); err != nil {
			return
		}
	}
	return
}
