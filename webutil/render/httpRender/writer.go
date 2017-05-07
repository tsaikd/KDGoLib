package httpRender

import (
	"io"
	"net/http"
)

func (t *renderImpl) GetIOWriter() io.Writer {
	return &ioWriter{t}
}

type ioWriter struct {
	*renderImpl
}

// Write is a low level API for io.Writer interface, use JSON, Error, WriteResponse in general case
func (t *ioWriter) Write(p []byte) (n int, err error) {
	if !t.written {
		// first Write call
		t.written = true
	}

	switch t.status {
	case 0:
		// first Write call
		t.status = http.StatusOK
		t.w.WriteHeader(http.StatusOK)
	case http.StatusNotModified:
		return 0, nil
	}

	writer := getWriter(t.renderImpl)
	defer func() {
		t.size += writer.size
	}()
	return writer.Write(p)
}
