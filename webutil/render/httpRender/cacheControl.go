package httpRender

import (
	"runtime/debug"
	"strconv"
	"time"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/webutil/render"
)

// http headers
const (
	HeaderCacheControl = "Cache-Control"
)

// errors
var (
	ErrorUnsupportedCacheControlType1 = errutil.NewFactory("unsupported cache-control type: %v")
)

func (t *renderImpl) SetCacheControl(ctype render.CacheControlType, maxAge time.Duration) {
	switch ctype {
	case render.CacheControlNoCache, render.CacheControlNoStore:
		t.w.Header().Set(HeaderCacheControl, ctype.String())
	case render.CacheControlPublic, render.CacheControlPrivate:
		value := ctype.String() + ", max-age=" + strconv.FormatInt(int64(maxAge/time.Second), 10)
		t.w.Header().Set(HeaderCacheControl, value)
	default:
		errutil.Trace(ErrorUnsupportedCacheControlType1.New(nil, ctype))
		debug.PrintStack()
	}
}
