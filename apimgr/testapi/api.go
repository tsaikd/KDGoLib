package testapi

import "github.com/tsaikd/KDGoLib/apimgr"

type TestAPI struct{}

var Manager = apimgr.NewManager(TestAPI{})
