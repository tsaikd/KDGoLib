package version

import (
	"encoding/json"
	"fmt"
)

type Version struct {
	VERSION   string `json:"version"`
	BUILDTIME string `json:"buildtime,omitempty"`
	GITCOMMIT string `json:"gitcommit,omitempty"`
}

var (
	VERSION   = "0.0.0"
	BUILDTIME string
	GITCOMMIT string
)

func Get() (ver Version) {
	ver = Version{
		VERSION:   VERSION,
		BUILDTIME: BUILDTIME,
		GITCOMMIT: GITCOMMIT,
	}
	return
}

func Json() (output string, err error) {
	var (
		raw []byte
	)
	ver := Get()
	if raw, err = json.MarshalIndent(ver, "", "\t"); err != nil {
		return
	}
	output = string(raw)
	return
}

func String() (output string) {
	output = VERSION
	if BUILDTIME != "" {
		output += fmt.Sprintf(" [%s]", BUILDTIME)
	}
	if GITCOMMIT != "" {
		output += fmt.Sprintf(" (%s)", GITCOMMIT)
	}
	return
}
