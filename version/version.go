package version

import (
	"encoding/json"
	"fmt"
)

// exported variables
var (
	VERSION   = "0.0.0"
	BUILDTIME string
	GITCOMMIT string
	GODEPS    string
)

// Version contains info for Get response
type Version struct {
	VERSION   string      `json:"version"`
	BUILDTIME string      `json:"buildtime,omitempty"`
	GITCOMMIT string      `json:"gitcommit,omitempty"`
	GODEPS    interface{} `json:"godeps,omitempty"`
}

// Get return Version info
func Get() (ver Version) {
	var godeps interface{}
	if GODEPS != "" {
		json.Unmarshal([]byte(GODEPS), &godeps)
	}
	return Version{
		VERSION:   VERSION,
		BUILDTIME: BUILDTIME,
		GITCOMMIT: GITCOMMIT,
		GODEPS:    godeps,
	}
}

// Json return version info with JSON format
func Json() (output string, err error) {
	var raw []byte
	ver := Get()
	if raw, err = json.MarshalIndent(ver, "", "\t"); err != nil {
		return
	}
	output = string(raw)
	return
}

// String return version info with string format
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
