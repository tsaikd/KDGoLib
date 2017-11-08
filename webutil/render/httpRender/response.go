package httpRender

import "net/http"

func (t *renderImpl) GetResponseHeader() http.Header {
	return t.w.Header()
}

func (t *renderImpl) GetResponseWriter() http.ResponseWriter {
	return &wrapResponseWriter{t}
}

type wrapResponseWriter struct {
	*renderImpl
}

func (t wrapResponseWriter) Header() http.Header {
	return t.w.Header()
}

func (t *wrapResponseWriter) Write(p []byte) (n int, err error) {
	return t.write(p)
}

func (t *wrapResponseWriter) WriteHeader(statusCode int) {
	t.writeStatus(statusCode)
}
