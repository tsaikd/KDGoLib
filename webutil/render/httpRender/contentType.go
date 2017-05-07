package httpRender

// http headers
const (
	HeaderContentType = "Content-Type"
)

func (t *renderImpl) SetContentType(ctype string) {
	t.w.Header().Set(HeaderContentType, ctype)
}
