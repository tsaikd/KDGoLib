package httpRender

import "net/http"

func (t *renderImpl) Redirect(status int, location string) {
	if !t.expectWritten(false) {
		return
	}
	t.written = true

	if !t.expectStatus(0) {
		return
	}
	t.status = status

	http.Redirect(t.w, t.req, location, status)
}
