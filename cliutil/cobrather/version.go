package cobrather

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/version"
)

// VersionModule provide module of version,
// Set nil to disable version command
var VersionModule = &Module{
	Use:   "version",
	Short: "Show version detail",
	RunE: func(cmd *cobra.Command, args []string) error {
		verjson, err := version.Json()
		if err != nil {
			return err
		}
		fmt.Println(verjson)
		return nil
	},
}
