package httpRender

import (
	"bytes"
	"net/http"
	"time"

	"github.com/tsaikd/KDGoLib/webutil/render"
)

// Render contains all interfaces can be used for render http response
type Render interface {
	render.Body
	render.CacheControl
	render.ContentType
	render.Cookie
	render.Error
	render.JSON
	render.LastModified
	render.Redirect
	render.Request
	render.Response
	render.Status
	render.Write
	render.Writer
}

// New return render instance
func New(w http.ResponseWriter, req *http.Request, options ...Option) Render {
	maxBufferSize := int64(1 << 20) // 1 MB
	errorPathTrimPrefixList := []string{}

	for _, option := range options {
		switch opt := option.(type) {
		case OptionMaxBufferSize:
			maxBufferSize = int64(opt)
		case OptionErrorPathTrimPrefix:
			errorPathTrimPrefixList = append(errorPathTrimPrefixList, string(opt))
		}
	}

	return &renderImpl{
		w:                       w,
		req:                     req,
		buffer:                  &bytes.Buffer{},
		maxBufferSize:           maxBufferSize,
		errorPathTrimPrefixList: errorPathTrimPrefixList,
	}
}

// Option for Render
type Option interface{}

// OptionMaxBufferSize max buffer size for render, default: 1 MB
type OptionMaxBufferSize int64

// OptionErrorPathTrimPrefix prefix will be trim in error path
type OptionErrorPathTrimPrefix string

type renderImpl struct {
	w                       http.ResponseWriter
	req                     *http.Request
	buffer                  *bytes.Buffer
	maxBufferSize           int64
	errorPathTrimPrefixList []string

	written      bool
	status       int
	size         int64
	err          error
	lastModified time.Time
}
