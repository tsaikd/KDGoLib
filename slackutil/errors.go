package slackutil

import "github.com/tsaikd/KDGoLib/errutil"

// errors
var (
	ErrSlackAuthTestFailed = errutil.NewFactory("slack client auth test failed")
	ErrSlackFormatFailed1  = errutil.NewFactory("format slack message failed %q")
)
