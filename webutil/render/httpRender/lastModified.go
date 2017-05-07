package httpRender

import (
	"net/http"
	"time"
)

// http headers
const (
	HeaderLastModified    = "Last-Modified"
	HeaderIfModifiedSince = "If-Modified-Since"
)

func (t *renderImpl) LastModified(ts time.Time) (notModified bool) {
	if !expectWritten(t, false) {
		return false
	}
	if t.lastModified.Before(ts) {
		t.lastModified = ts
		t.w.Header().Set(HeaderLastModified, t.lastModified.Format(http.TimeFormat))
	}
	if matchLastModified(t.req.Header, t.lastModified) {
		t.written = true
		t.status = http.StatusNotModified
		t.w.WriteHeader(http.StatusNotModified)
		return true
	}
	return false
}

func matchLastModified(reqHeader http.Header, ts time.Time) bool {
	if ts.IsZero() {
		return false
	}
	return reqHeader.Get(HeaderIfModifiedSince) == ts.Format(http.TimeFormat)
}
