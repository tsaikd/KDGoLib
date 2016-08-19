package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/kardianos/osext"
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
	Version   string      `json:"version"`
	GoVersion string      `json:"goversion"`
	BuildTime string      `json:"buildtime,omitempty"`
	GitCommit string      `json:"gitcommit,omitempty"`
	Godeps    interface{} `json:"godeps,omitempty"`
}

// Get return Version info
func Get() (ver Version) {
	var godeps interface{}
	if GODEPS != "" {
		_ = json.Unmarshal([]byte(GODEPS), &godeps)
	}
	if BUILDTIME == "" {
		if exectime, err := getExecModifyTime(); err == nil {
			BUILDTIME = exectime.Format(time.RFC1123)
		}
	}
	return Version{
		Version:   VERSION,
		GoVersion: runtime.Version(),
		BuildTime: BUILDTIME,
		GitCommit: GITCOMMIT,
		Godeps:    godeps,
	}
}

// JSON return version info with JSON format
func JSON() (output string, err error) {
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

func getExecModifyTime() (modtime time.Time, err error) {
	execFileName, err := osext.Executable()
	if err != nil {
		return
	}

	fi, err := os.Stat(execFileName)
	if err != nil {
		return
	}

	return fi.ModTime(), nil
}
