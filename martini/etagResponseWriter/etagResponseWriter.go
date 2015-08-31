package etagResponseWriter

import (
	"net/http"

	"github.com/go-martini/martini"
)

func ETagResponseWriter(config *ETagConfig) martini.Handler {
	if config == nil {
		config = NewETagConfig()
	}
	return func(c martini.Context, res http.ResponseWriter, req *http.Request) {
		resp := &etagResponseWriter{
			config: *config,
			rw:     res,
			req:    req,
		}
		if resp.isEnableMethod(req.Method) && !resp.isIgnoreIfHeader(req) {
			c.MapTo(resp, (*http.ResponseWriter)(nil))
		}
		c.Next()
	}
}

type etagResponseWriter struct {
	config      ETagConfig
	rw          http.ResponseWriter
	req         *http.Request
	status      int
	writeStatus bool
}

func (t *etagResponseWriter) Header() http.Header {
	return t.rw.Header()
}

func (t *etagResponseWriter) Write(b []byte) (int, error) {
	if len(b) < t.config.MinBodyLength {
		return t.write(b)
	}

	if !t.isEnableStatus(t.status) {
		return t.write(b)
	}

	etag := t.config.HashFunc(b)

	if match, ok := t.req.Header["If-None-Match"]; ok && match[0] == etag {
		t.WriteHeader(http.StatusNotModified)
		b = []byte{}
	} else {
		t.Header().Set("ETag", etag)
	}

	return t.write(b)
}

func (t *etagResponseWriter) WriteHeader(s int) {
	t.status = s
}

func (t *etagResponseWriter) isEnableMethod(method string) bool {
	if enable, ok := t.config.EnableMethod[method]; ok {
		return enable
	}
	return false
}

func (t *etagResponseWriter) isEnableStatus(status int) bool {
	if enable, ok := t.config.EnableStatus[status]; ok {
		return enable
	}
	return false
}

func (t *etagResponseWriter) isIgnoreIfHeader(req *http.Request) bool {
	for header, value := range t.config.IgnoreIfHeader {
		if req.Header.Get(header) == value {
			return true
		}
	}
	return false
}

func (t *etagResponseWriter) write(b []byte) (int, error) {
	if !t.writeStatus && t.status > 0 {
		t.rw.WriteHeader(t.status)
		t.writeStatus = true
	}

	return t.rw.Write(b)
}
