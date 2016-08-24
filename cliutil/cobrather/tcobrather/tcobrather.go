package tcobrather

import (
	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
)

// NewTest create TestModule instance
func NewTest(module *cobrather.Module) *TestModule {
	return &TestModule{
		command: module.MustNewRootCommand(nil),
	}
}

// TestModule used for testing cmder.Module
type TestModule struct {
	command *cobra.Command
}

// Setup run before, action in module
func (t TestModule) Setup() (err error) {
	if t.command.PreRunE != nil {
		if err = t.command.PreRunE(t.command, []string{}); err != nil {
			return
		}
	}

	if t.command.RunE != nil {
		if err = t.command.RunE(t.command, []string{}); err != nil {
			return
		}
	}

	return
}

// Teardown run after in module
func (t TestModule) Teardown() (err error) {
	if t.command.PostRunE != nil {
		if err = t.command.PostRunE(t.command, []string{}); err != nil {
			return
		}
	}
	return
}
