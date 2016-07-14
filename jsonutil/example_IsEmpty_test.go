package jsonutil_test

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/jsonutil"
)

type emptyStruct struct {
	String string
	empty  bool
}

func (t emptyStruct) IsEmpty() bool {
	return t.empty
}

// Custom struct with IsEmpty() interface will change the marshal output
func ExampleIsEmpty() {
	data1, _ := jsonutil.MarshalJSON(emptyStruct{
		String: "text",
		empty:  false,
	})
	data2, _ := jsonutil.MarshalJSON(emptyStruct{
		String: "text",
		empty:  true,
	})
	fmt.Printf("%s\n", data1)
	fmt.Printf("%s\n", data2)
	// Output:
	// {"String":"text"}
	// null
}
