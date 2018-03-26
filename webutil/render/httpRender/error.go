package httpRender

import (
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/tsaikd/KDGoLib/errutil"
)

func (t *renderImpl) Error(err error) {
	if err == nil {
		return
	}

	t.err = err

	errjson, err := errutil.NewJSON(errutil.NewErrorsSkip(1, err))
	errutil.Trace(err)

	output := struct {
		Status int `json:"status,omitempty"`
		*errutil.ErrorJSON
	}{
		Status:    http.StatusNotFound,
		ErrorJSON: errjson,
	}

	for _, trim := range t.errorPathTrimPrefixList {
		trimpath := strings.TrimPrefix(output.ErrorPath, trim)
		if trimpath != output.ErrorPath { // only match one trim rule
			output.ErrorPath = trimpath
			break
		}
	}

	if operr, ok := t.err.(*net.OpError); ok {
		if syserr, ok := operr.Err.(*os.SyscallError); ok {
			if syserr.Err == syscall.EPIPE {
				// do not send error message if the error cause from write a broken pipe
				return
			}
		}
	}

	t.WriteResponse(nil, output.Status, output)
}

func (t renderImpl) GetError() error {
	return t.err
}
