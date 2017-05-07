package httpRender

import "net/http"

func (t *renderImpl) SetCookie(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	if !expectWritten(t, false) {
		return
	}
	http.SetCookie(t.w, cookie)
}
