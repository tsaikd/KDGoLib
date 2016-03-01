package cmdutil

import (
	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	commands = map[string]cli.Command{}
)

// errors
var (
	ErrorCommandDefined1 = errutil.NewFactory("command %s defined")
)

func AddCommand(command cli.Command) cli.Command {
	if _, ok := commands[command.Name]; ok {
		panic(ErrorCommandDefined1.New(nil, command.Name))
	} else {
		commands[command.Name] = command
	}
	return command
}

func AllCommands() (retcommands []cli.Command) {
	for _, cmd := range commands {
		retcommands = append(retcommands, cmd)
	}
	return
}
