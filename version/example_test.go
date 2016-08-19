package version_test

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/version"
)

func init() {
	version.VERSION = "0.0.1"
}

func Example() {
	fmt.Println(version.String())

	ver := version.Get()
	fmt.Println(ver.Version)

	// Output:
	// 0.0.1
	// 0.0.1
}
