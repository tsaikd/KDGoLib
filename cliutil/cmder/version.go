package cmder

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/version"
)

// VersionModule provide module of version,
// Set nil to disable version command
var VersionModule = NewModule("version").
	SetUsage("Show version detail").
	SetAction(versionAction)

func versionAction(c *cli.Context) (err error) {
	verjson, err := version.Json()
	if err != nil {
		return
	}
	fmt.Println(verjson)
	return
}
