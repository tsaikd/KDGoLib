package cmder

import (
	"github.com/tsaikd/KDGoLib/errutil"
	"gopkg.in/urfave/cli.v2"
)

// WrapAction wrap action with default error handler
func WrapAction(action cli.ActionFunc) cli.ActionFunc {
	if action == nil {
		return nil
	}

	return func(c *cli.Context) (err error) {
		if err = action(c); err != nil {
			switch err.(type) {
			case cli.ExitCoder:
				return
			default:
				var message string
				formatter := errutil.NewConsoleFormatter("; ")
				if message, err = formatter.FormatSkip(err, 1); err != nil {
					panic(err)
				}
				return cli.Exit(message, 1)
			}
		}
		return
	}
}

// WrapMainAction wrap main action with default error handler and command helper
func WrapMainAction(action cli.ActionFunc) cli.ActionFunc {
	if action == nil {
		return nil
	}

	return func(c *cli.Context) (err error) {
		args := c.Args()
		if args.Present() {
			return cli.ShowCommandHelp(c, args.First())
		}

		return WrapAction(action)(c)
	}
}

// Actions slice of cli.ActionFunc
type Actions []cli.ActionFunc

// Add action
func (t *Actions) Add(actions ...cli.ActionFunc) {
	for _, action := range actions {
		if action == nil {
			continue
		}
		*t = append(*t, action)
	}
}

// Wrap return wrapped function for running all actions with WrapAction
func (t Actions) Wrap() cli.ActionFunc {
	if len(t) < 1 {
		return nil
	}

	return WrapAction(func(c *cli.Context) (err error) {
		for _, action := range t {
			if err = action(c); err != nil {
				return
			}
		}
		return nil
	})
}

// WrapMain return wrapped function for running all actions with WrapMainAction
func (t Actions) WrapMain(mainAction cli.ActionFunc) cli.ActionFunc {
	return WrapMainAction(func(c *cli.Context) (err error) {
		for _, action := range t {
			if err = action(c); err != nil {
				return
			}
		}

		if mainAction != nil {
			return mainAction(c)
		}

		args := c.Args()
		if args.Present() {
			return cli.ShowCommandHelp(c, args.First())
		}

		cli.ShowAppHelp(c)
		return nil
	})
}
