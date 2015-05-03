package flagutil

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var (
	flags = map[string]cli.Flag{}
)

func AddBoolFlag(cliflag cli.BoolFlag) cli.BoolFlag {
	var (
		ok bool
	)
	if _, ok = flags[cliflag.Name]; ok {
		panic(fmt.Errorf("flag %s defined", cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddIntFlag(cliflag cli.IntFlag) cli.IntFlag {
	var (
		ok bool
	)
	if _, ok = flags[cliflag.Name]; ok {
		panic(fmt.Errorf("flag %s defined", cliflag.Name))
	} else {
		flags[cliflag.Name] = cliflag
	}
	return cliflag
}

func AddStringFlag(cliflag cli.StringFlag) cli.StringFlag {
	var (
		ok bool
	)
	if _, ok = flags[cliflag.Name]; ok {
		panic(fmt.Errorf("flag %s defined", cliflag.Name))
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
