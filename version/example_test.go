package version_test

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/version"
)

func init() {
	version.VERSION = "0.0.1"
}

func ExampleShowVersionString() {
	fmt.Println(version.String())
}

func ExampleGet() {
	ver := version.Get()
	fmt.Println(ver.Version)

	// Output:
	// 0.0.1
}
