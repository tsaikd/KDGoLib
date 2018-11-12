package slackutil

import (
	"bytes"
	"strings"

	"github.com/nlopes/slack"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/runtimecaller"
)

// Formatter used to format error object to slack message
type Formatter struct {
	// show error position with short filename
	ShortFile bool
	// hide error position with line number
	HideLine bool
	// show customized user name for slack message
	UserName string
	// filter runtime caller
	CallerFilter []runtimecaller.Filter
	// replace package name for securify
	ReplacePackages map[string]string
}

func (t *Formatter) formatSkip(text string, skip int, stack bool) (msgopts []slack.MsgOption, err error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	// main message
	msgopts = append(msgopts, slack.MsgOptionText(text, true))

	if t.UserName != "" {
		param := slack.NewPostMessageParameters()
		param.Username = t.UserName
		msgopts = append(msgopts, slack.MsgOptionPostMessageParameters(param))
	}

	// attachment for full call stack info
	if stack {
		callstack := &bytes.Buffer{}
		callinfos := runtimecaller.ListByFilters(skip+1, t.CallerFilter...)
		for _, callinfo := range callinfos {
			if callstack.Len() > 0 {
				if _, err = callstack.WriteString("\n"); err != nil {
					return
				}
			}
			if _, err = errutil.WriteCallInfo(callstack, callinfo, !t.ShortFile, !t.HideLine, t.ReplacePackages); err != nil {
				return
			}
		}

		msgopts = append(msgopts, slack.MsgOptionAttachments(slack.Attachment{
			Text: callstack.String(),
		}))
	}

	return
}
