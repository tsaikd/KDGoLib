package climgr

import (
	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/logutil"
)

// errors
var (
	ErrorCommandRegisted1 = errutil.NewFactory("Command already registed: %q")
)

// AddCommand add command to manager
func (t *Manager) AddCommand(command cli.Command) (err error) {
	key := command.Name
	if _, exist := t.commandNameMap[key]; exist {
		return ErrorCommandRegisted1.New(nil, key)
	}
	t.commandNameMap[key] = command
	return
}

// DeleteCommand delete command from manager
func (t *Manager) DeleteCommand(command cli.Command) {
	key := command.Name
	delete(t.commandNameMap, key)
}

// GetCommands return registed commands
func (t *Manager) GetCommands() (results []cli.Command) {
	for _, command := range t.commandNameMap {
		results = append(results, command)
	}
	return
}

// WrapAction wrap action, handle return error with errutil.Trace()
func WrapAction(action func(context *cli.Context) error) func(context *cli.Context) {
	return func(context *cli.Context) {
		if err := action(context); err != nil {
			errutil.Trace(err)
		}
	}
}

// WrapActionLogger wrap action, handle return error with logger
func WrapActionLogger(action func(context *cli.Context) error, logger logutil.StdLogger) func(context *cli.Context) {
	return func(context *cli.Context) {
		if err := action(context); err != nil {
			logger.Fatalln(err)
		}
	}
}
