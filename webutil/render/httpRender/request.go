package httpRender

import "net/http"

func (t *renderImpl) GetRequest() *http.Request {
	return t.req
}
