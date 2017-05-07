package httpRender

import (
	"bytes"
	"io"
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
	if !expectWritten(t, false) {
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

	if !expectStatus(t, 0) {
		return
	}
	t.status = status
	t.w.WriteHeader(status)
	if !bodyAllowedForStatus(status) && data == nil {
		return
	}

	writer := getWriter(t)
	defer func() {
		t.size += writer.size
	}()
	switch value := data.(type) {
	case []byte:
		if len(value) < 1 {
			return
		}
		if _, err := writer.Write(value); err != nil {
			errutil.Trace(err)
		}
	default:
		jsonenc := jsonex.NewEncoder(writer)
		errutil.Trace(jsonenc.Encode(value))
	}
}

func (t renderImpl) IsWritten() bool {
	return t.written
}

func expectWritten(r *renderImpl, written bool) bool {
	if written == r.written {
		return true
	}

	errutil.Trace(ErrorUnexpectedWrittenStatus2.New(nil, written, r.written))
	debug.PrintStack()
	return false
}

func expectStatus(r *renderImpl, status int) bool {
	if status == r.status {
		return true
	}

	errutil.Trace(ErrorUnexpectedStatusCode2.New(nil, status, r.status))
	debug.PrintStack()
	return false
}

type writerLogSize struct {
	writer io.Writer
	size   int64

	buffer        *bytes.Buffer
	maxBufferSize int64
}

func (t *writerLogSize) Write(p []byte) (n int, err error) {
	n, err = t.writer.Write(p)
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

var _ io.Writer = &writerLogSize{}

// getWriter return io.Writer with BufferResponse flag
func getWriter(r *renderImpl) (writer *writerLogSize) {
	return &writerLogSize{
		writer: r.w,

		buffer:        r.buffer,
		maxBufferSize: r.maxBufferSize,
	}
}
