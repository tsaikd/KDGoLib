package testmethod11

import (
	"github.com/tsaikd/KDGoLib/apimgr"
	"github.com/tsaikd/KDGoLib/apimgr/testapi"
)

func init() {
	testapi.Manager.Add(
		apimgr.Definition{
			Description: `
				Test api 1 method 11
			`,
			Method:  "GET",
			Pattern: "/1/testmethod11/testres1",
			Request: TestAPI11{},
		},
	)
}

type TestAPI11 struct{}
