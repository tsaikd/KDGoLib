package httpRender

import (
	"runtime/debug"

	"github.com/tsaikd/KDGoLib/errutil"
)

func (t renderImpl) expectStatus(status int) bool {
	if status == t.status {
		return true
	}

	errutil.Trace(ErrorUnexpectedStatusCode2.New(nil, status, t.status))
	debug.PrintStack()
	return false
}

func (t *renderImpl) setStatus(status int) {
	if !t.expectStatus(0) {
		return
	}
	t.status = status
}

func (t *renderImpl) writeStatus(status int) {
	if !t.expectStatus(0) {
		return
	}
	t.status = status
	t.w.WriteHeader(status)
}

func (t renderImpl) GetStatus() int {
	return t.status
}
