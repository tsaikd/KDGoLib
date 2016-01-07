package climgr

import (
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func Test_command(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	manager := NewManager()

	err := manager.AddCommand(cli.Command{
		Name: "A",
	})
	assert.NoError(err)

	commands := manager.GetCommands()
	assert.Len(commands, 1)
	assert.Equal("A", commands[0].Name)

	err = manager.AddCommand(cli.Command{
		Name: "B",
	})
	assert.NoError(err)

	commands = manager.GetCommands()
	assert.Len(commands, 2)

	err = manager.AddCommand(cli.Command{
		Name: "C",
	})
	assert.NoError(err)

	commands = manager.GetCommands()
	assert.Len(commands, 3)
}
