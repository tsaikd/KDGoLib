package jsonutil_test

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/jsonutil"
)

func ExampleMarshalJSON() {
	type myStruct struct {
		String string
	}
	data, _ := jsonutil.MarshalJSON(myStruct{
		String: "text",
	})
	fmt.Printf("%s\n", data)
	// Output:
	// {"String":"text"}
}

func ExampleMarshalJSON_omitdefault_string() {
	type myStruct struct {
		String string `json:",omitdefault" default:"text"`
	}
	data, _ := jsonutil.MarshalJSON(myStruct{
		String: "text",
	})
	fmt.Printf("%s\n", data)
	// Output:
	// {}
}

func ExampleMarshalJSON_omitempty_struct() {
	type subStruct struct {
		String string `json:",omitdefault" default:"text"`
	}
	type myStruct struct {
		SubStruct subStruct `json:",omitempty"`
	}
	data, _ := jsonutil.MarshalJSON(myStruct{
		subStruct{
			String: "text",
		},
	})
	fmt.Printf("%s\n", data)
	// Output:
	// {}
}
