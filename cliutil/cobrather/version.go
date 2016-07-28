package cobrather

import (
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
	FlagNumber = &BoolFlag{
		Name:      "number",
		ShortHand: "n",
		Usage:     "only show the version number",
	}
	FlagContains = &StringFlag{
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
		FlagNumber,
		FlagContains,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if FlagNumber.Bool() {
			fmt.Println(version.VERSION)
			return nil
		}
		if rangeStr := FlagContains.String(); rangeStr != "" {
			if err := checkVersionRange(rangeStr); err != nil {
				return err
			}
			return nil
		}

		verjson, err := version.JSON()
		if err != nil {
			return err
		}
		fmt.Println(verjson)
		return nil
	},
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
