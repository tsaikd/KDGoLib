package cobrather

import (
	"context"
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/version"
)

// erors
var (
	ErrorVersionNotInRange2 = errutil.NewFactory("current version %q not in range %q")
	ErrorInvalidRange1      = errutil.NewFactory("invalid range string %q")
)

// command line flags
var (
	flagNumber = &BoolFlag{
		Name:      "number",
		ShortHand: "n",
		Usage:     "only show the version number",
	}
	flagContains = &StringFlag{
		Name:      "contains",
		ShortHand: "c",
		Usage:     "version is inside the range",
	}
)

// VersionModule provide module of version,
// Set nil to disable version command
var VersionModule = &Module{
	Use:   "version",
	Short: "Show version detail",
	Example: strings.TrimSpace(`
version -n
version -c ">=0.3.5"
version -c ">=0.3.5 , <1"
	`),
	Flags: []Flag{
		flagNumber,
		flagContains,
	},
	RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
		return showVersion(flagNumber.Bool(), flagContains.String())
	},
}

func showVersion(onlyNumber bool, versionRange string) (err error) {
	if onlyNumber {
		fmt.Println(version.VERSION)
		return nil
	}

	if versionRange != "" {
		return checkVersionRange(versionRange)
	}

	verjson, err := version.JSON()
	if err != nil {
		return err
	}
	fmt.Println(verjson)

	return nil
}

func checkVersionRange(rangeStr string) (err error) {
	ver, err := semver.NewVersion(version.VERSION)
	if err != nil {
		return
	}

	verRange, err := semver.NewConstraint(rangeStr)
	if err != nil {
		return ErrorInvalidRange1.New(err, rangeStr)
	}

	if !verRange.Check(ver) {
		return ErrorVersionNotInRange2.New(nil, version.VERSION, rangeStr)
	}
	return nil
}
