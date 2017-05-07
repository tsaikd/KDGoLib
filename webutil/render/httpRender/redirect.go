package httpRender

import "net/http"

func (t *renderImpl) Redirect(status int, location string) {
	if !expectWritten(t, false) {
		return
	}
	t.written = true

	if !expectStatus(t, 0) {
		return
	}
	t.status = status

	http.Redirect(t.w, t.req, location, status)
}
