package flagutil

import (
	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	flags = map[string]cli.Flag{}
)

// errors
var (
	ErrorFlagDefined1 = errutil.NewFactory("flag %s defined")
)

func AddBoolFlag(cliflag cli.BoolFlag) cli.BoolFlag {
	if _, ok := flags[cliflag.Name]; ok {
		panic(ErrorFlagDefined1.New(nil, cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddIntFlag(cliflag cli.IntFlag) cli.IntFlag {
	if _, ok := flags[cliflag.Name]; ok {
		panic(ErrorFlagDefined1.New(nil, cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddIntSliceFlag(cliflag cli.IntSliceFlag) cli.IntSliceFlag {
	if _, ok := flags[cliflag.Name]; ok {
		panic(ErrorFlagDefined1.New(nil, cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddStringFlag(cliflag cli.StringFlag) cli.StringFlag {
	if _, ok := flags[cliflag.Name]; ok {
		panic(ErrorFlagDefined1.New(nil, cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddStringSliceFlag(cliflag cli.StringSliceFlag) cli.StringSliceFlag {
	if _, ok := flags[cliflag.Name]; ok {
		panic(ErrorFlagDefined1.New(nil, cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AllFlags() (retflags []cli.Flag) {
	for _, f := range flags {
		retflags = append(retflags, f)
	}
	return
}
