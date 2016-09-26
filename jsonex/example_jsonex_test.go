package jsonex_test

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/jsonex"
)

// omitdefault flag with default tag on struct field
func ExampleMarshal_omitdefault() {
	type myStruct struct {
		Bool        bool        `json:",omitdefault" default:"true"`
		Int         int64       `json:",omitdefault" default:"-9527"`
		Uint        uint64      `json:",omitdefault" default:"9527"`
		Float       float64     `json:",omitdefault" default:"3.14"`
		String      string      `json:",omitdefault" default:"text"`
		StringPtr   *string     `json:",omitdefault" default:"textptr"`
		StringSlice []string    `json:",omitdefault"`
		Interface   interface{} `json:",omitdefault"`
	}
	textptr := "textptr"
	data, err := jsonex.Marshal(myStruct{
		Bool:      true,
		Int:       -9527,
		Uint:      9527,
		Float:     3.14,
		String:    "text",
		StringPtr: &textptr,
	})
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {}
}

type marshalStruct struct {
	MarshalText string `json:",omitempty"`
}

func (t marshalStruct) MarshalJSON() ([]byte, error) {
	return jsonex.Marshal(map[string]interface{}{})
}

// omitempty flag on struct field
func ExampleMarshal_omitemptyOnStruct() {
	type subStruct struct {
		String string `json:",omitdefault" default:"text"`
		Int    int64
	}
	type PubStruct struct {
		String string
		Ignore string `json:"-"`
	}
	type WrapPubStruct struct {
		PubStruct
	}
	type myStruct struct {
		subStruct
		WrapPubStruct WrapPubStruct `json:",omitempty"`
		MarshalStruct marshalStruct `json:",omitempty"`
		Float         float64
	}
	data, err := jsonex.Marshal(myStruct{
		subStruct: subStruct{
			String: "text",
		},
	})
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"Int":0,"Float":0}
}
