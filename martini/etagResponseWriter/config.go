package etagResponseWriter

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

// ETagHashFunc hash function for generate etag from response body
type ETagHashFunc func(data []byte) string

// ETagConfig config for ETag
type ETagConfig struct {
	MinBodyLength       int
	EnableMethod        map[string]bool
	EnableStatus        map[int]bool
	IgnoreIfHeaderExist map[string]bool
	IgnoreIfHeaderValue map[string]string
	HashFunc            ETagHashFunc
}

// NewETagConfig return default ETagConfig
func NewETagConfig() *ETagConfig {
	return &ETagConfig{
		MinBodyLength: 1024,
		EnableMethod: map[string]bool{
			"GET": true,
		},
		EnableStatus: map[int]bool{
			0:             true,
			http.StatusOK: true,
		},
		IgnoreIfHeaderExist: map[string]bool{
			"If-Modified-Since": true,
		},
		IgnoreIfHeaderValue: map[string]string{
			"Upgrade": "websocket",
		},
		HashFunc: func(data []byte) string {
			hash := md5.Sum(data)
			return fmt.Sprintf("%x", hash)
		},
	}
}

// SetMinBodyLength if response body size less than length, do not use etag
func (t *ETagConfig) SetMinBodyLength(length int) *ETagConfig {
	t.MinBodyLength = length
	return t
}

// AddMethod allow request method for etag
func (t *ETagConfig) AddMethod(method string) *ETagConfig {
	t.EnableMethod[method] = true
	return t
}

// AddStatus allow response status for etag
func (t *ETagConfig) AddStatus(status int) *ETagConfig {
	t.EnableStatus[status] = true
	return t
}

// AddIgnoreHeaderExist ignore etag if header exist
func (t *ETagConfig) AddIgnoreHeaderExist(header string, exist bool) *ETagConfig {
	t.IgnoreIfHeaderExist[header] = exist
	return t
}

// AddIgnoreHeaderValue ignore etag if header equal value
func (t *ETagConfig) AddIgnoreHeaderValue(header string, value string) *ETagConfig {
	t.IgnoreIfHeaderValue[header] = value
	return t
}

// SetHashFunc set etag hash function
func (t *ETagConfig) SetHashFunc(hashFunc ETagHashFunc) *ETagConfig {
	t.HashFunc = hashFunc
	return t
}
