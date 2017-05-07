package httpRender

import "net/http"

func (t *renderImpl) GetResponseHeader() http.Header {
	return t.w.Header()
}
