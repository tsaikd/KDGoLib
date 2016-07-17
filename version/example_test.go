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
	verjson, err := version.Json()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(verjson)
	// Output:
	// 0.0.1
	// {
	//	"version": "0.0.1"
	// }
}
