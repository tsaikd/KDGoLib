package httpRender

import (
	"net/http"
	"runtime/debug"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/jsonex"
)

// errors
var (
	ErrorUnexpectedWrittenStatus2 = errutil.NewFactory("expect written status is %v but got %v")
	ErrorUnexpectedStatusCode2    = errutil.NewFactory("expect status code is %v but got %v")
)

const (
	contentJSON = "application/json"
	charset     = "; charset=UTF-8"
)

func (t *renderImpl) WriteResponse(header http.Header, status int, data interface{}) {
	if !t.expectWritten(false) {
		return
	}
	t.written = true

	for headerKey := range header {
		t.w.Header().Set(headerKey, header.Get(headerKey))
	}

	switch data.(type) {
	case []byte:
	default:
		t.SetContentType(contentJSON + charset)
	}

	t.writeStatus(status)
	if !bodyAllowedForStatus(status) && data == nil {
		return
	}

	switch value := data.(type) {
	case []byte:
		if len(value) < 1 {
			return
		}
		if _, err := t.write(value); err != nil {
			errutil.Trace(err)
		}
	default:
		jsonenc := jsonex.NewEncoder(t.GetResponseWriter())
		errutil.Trace(jsonenc.Encode(value))
	}
}

func (t renderImpl) IsWritten() bool {
	return t.written
}

func (t renderImpl) expectWritten(written bool) bool {
	writeStatus := t.IsWritten()
	if written == writeStatus {
		return true
	}

	errutil.Trace(ErrorUnexpectedWrittenStatus2.New(nil, written, writeStatus))
	debug.PrintStack()
	return false
}

func (t *renderImpl) write(p []byte) (n int, err error) {
	if !t.written {
		t.written = true
	}

	switch t.status {
	case 0:
		// first Write call
		t.writeStatus(http.StatusOK)
	case http.StatusNotModified:
		return 0, nil
	}

	n, err = t.w.Write(p)
	t.size += int64(n)

	rest := t.maxBufferSize - int64(t.buffer.Len())
	if rest > 0 {
		n2 := int64(n)
		if rest < n2 {
			n2 = rest
		}
		t.buffer.Write(p[0:n2])
	}

	return
}
