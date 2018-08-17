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
	NAME      string
	VERSION   = "0.0.0"
	BUILDTIME string
	GITCOMMIT string
	GODEPS    string
)

// Version contains info for Get response
type Version struct {
	Name      string      `json:"name,omitempty"`
	Version   string      `json:"version"`
	GoVersion string      `json:"goversion"`
	BuildTime string      `json:"buildtime,omitempty"`
	GitCommit string      `json:"gitcommit,omitempty"`
	Godeps    interface{} `json:"godeps,omitempty"`
}

func (t Version) String() string {
	output := t.Name
	if output != "" {
		output += " "
	}
	output += t.Version
	if t.BuildTime != "" {
		output += fmt.Sprintf(" [%s]", t.BuildTime)
	}
	if t.GitCommit != "" {
		output += fmt.Sprintf(" (%s)", t.GitCommit)
	}
	return output
}

// JSON return version info with JSON format
func (t Version) JSON() (output string, err error) {
	var raw []byte
	ver := Get()
	if raw, err = json.MarshalIndent(ver, "", "\t"); err != nil {
		return
	}
	return string(raw), nil
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
		Name:      NAME,
		Version:   VERSION,
		GoVersion: runtime.Version(),
		BuildTime: BUILDTIME,
		GitCommit: GITCOMMIT,
		Godeps:    godeps,
	}
}

// JSON return version info with JSON format
func JSON() (output string, err error) {
	return Get().JSON()
}

// String return version info with string format
func String() (output string) {
	return Get().String()
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
