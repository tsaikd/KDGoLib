package render

import (
	"io"
	"net/http"
	"time"
)

// Body interface for get http response body instance
type Body interface {
	// GetBody return response body data
	GetBody() []byte
	// GetSize return response body size, not equal to len(GetBody()) because GetBody() exist only if BufferResponse flag set
	GetSize() int64
}

// CacheControl interface for set response cache-control header
type CacheControl interface {
	// SetCacheControl set response cache-control header
	SetCacheControl(ctype CacheControlType, maxAge time.Duration)
}

// ContentType interface for set response content-type header
type ContentType interface {
	// SetContentType set response content-type header
	SetContentType(ctype string)
}

// Cookie interface for handling http cookie to response
type Cookie interface {
	// SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
	SetCookie(cookie *http.Cookie)
}

// Error interface for render Error to response
type Error interface {
	// Error write error to response
	Error(err error)
	// GetError return error for the last Error() called
	GetError() error
}

// JSON interface for render JSON object to response
type JSON interface {
	// JSON render object to response in JSON format
	JSON(obj interface{})
}

// LastModified interface handle http last modified feature
type LastModified interface {
	// LastModified set last-modified header in response and check if-modified-since header in request, return true if header match
	LastModified(ts time.Time) (notModified bool)
}

// Redirect http request to location with status code
type Redirect interface {
	// Redirect http request to location with status code
	Redirect(status int, location string)
}

// Request interface for get http request instance
type Request interface {
	// GetRequest return http request instance
	GetRequest() *http.Request
}

// Response interface for get http response instance
type Response interface {
	// GetResponseHeader return http response header
	GetResponseHeader() http.Header
}

// Status interface for get status from Render which already set before
type Status interface {
	// GetStatus return response status code
	GetStatus() int
}

// Write interface for Render custom data, low level API
type Write interface {
	// WriteResponse write data to response, use JSON/Error instead if possible
	WriteResponse(header http.Header, status int, data interface{})
	// IsWritten return true if already write something
	IsWritten() bool
}

// Writer interface for get io.Writer API, low level API
type Writer interface {
	// GetIOWriter return io.Writer interface
	GetIOWriter() io.Writer
}
