package cmder

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/version"
)

// Name of application
var Name = "cmder"

// Usage of application
var Usage string

// Action application main action
var Action = func(c *cli.Context) (err error) {
	return
}

// Flags list all global flags for application
var Flags = []cli.Flag{}

// Commands list all commands for application
var Commands = []cli.Command{}

// DisableVersionCommand do not add version command if DisableVersionCommand == true
var DisableVersionCommand = false

// Main entry point
func Main() {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = Usage
	app.Version = version.String()
	app.Action = Action
	app.Flags = Flags

	if !DisableVersionCommand {
		Commands = append(Commands, cli.Command{
			Name:   "version",
			Usage:  "Show version detail",
			Action: WrapAction(versionAction),
		})
	}
	app.Commands = Commands

	errutil.Trace(app.Run(os.Args))
}

// WrapAction wrap action with default error handler
func WrapAction(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) (err error) {
		if err = action(c); err != nil {
			switch err.(type) {
			case *cli.ExitError:
				return
			default:
				var message string
				formatter := errutil.NewConsoleFormatter("; ")
				if message, err = formatter.FormatSkip(err, 1); err != nil {
					panic(err)
				}
				return cli.NewExitError(message, 1)
			}
		}
		return
	}
}

func versionAction(c *cli.Context) (err error) {
	verjson, err := version.Json()
	if err != nil {
		return
	}
	fmt.Println(verjson)
	return
}
