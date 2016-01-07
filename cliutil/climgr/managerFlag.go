package climgr

import (
	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorFlagRegisted1 = errutil.ErrorFactory("Flag already registed: %q")
)

// MustAddBoolFlag add bool flag to manager, panic if error
func (t *Manager) MustAddBoolFlag(flag cli.BoolFlag) cli.BoolFlag {
	if _, ok := t.flagNameMap[flag.Name]; ok {
		panic(ErrorFlagRegisted1.New(nil, flag.Name))
	} else {
		t.flagNameMap[flag.Name] = flag
	}
	return flag
}

// MustAddIntFlag add int flag to manager, panic if error
func (t *Manager) MustAddIntFlag(flag cli.IntFlag) cli.IntFlag {
	if _, ok := t.flagNameMap[flag.Name]; ok {
		panic(ErrorFlagRegisted1.New(nil, flag.Name))
	} else {
		t.flagNameMap[flag.Name] = flag
	}
	return flag
}

// MustAddIntSliceFlag add int slice flag to manager, panic if error
func (t *Manager) MustAddIntSliceFlag(flag cli.IntSliceFlag) cli.IntSliceFlag {
	if _, ok := t.flagNameMap[flag.Name]; ok {
		panic(ErrorFlagRegisted1.New(nil, flag.Name))
	} else {
		t.flagNameMap[flag.Name] = flag
	}
	return flag
}

// MustAddStringFlag add string flag to manager, panic if error
func (t *Manager) MustAddStringFlag(flag cli.StringFlag) cli.StringFlag {
	if _, ok := t.flagNameMap[flag.Name]; ok {
		panic(ErrorFlagRegisted1.New(nil, flag.Name))
	} else {
		t.flagNameMap[flag.Name] = flag
	}
	return flag
}

// MustAddStringSliceFlag add string slice flag to manager, panic if error
func (t *Manager) MustAddStringSliceFlag(flag cli.StringSliceFlag) cli.StringSliceFlag {
	if _, ok := t.flagNameMap[flag.Name]; ok {
		panic(ErrorFlagRegisted1.New(nil, flag.Name))
	} else {
		t.flagNameMap[flag.Name] = flag
	}
	return flag
}

// GetFlags return registed flags
func (t *Manager) GetFlags() (results []cli.Flag) {
	for _, flag := range t.flagNameMap {
		results = append(results, flag)
	}
	return
}
