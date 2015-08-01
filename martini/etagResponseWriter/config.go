package etagResponseWriter

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

type ETagHashFunc func(data []byte) string

type ETagConfig struct {
	MinBodyLength  int
	EnableMethod   map[string]bool
	EnableStatus   map[int]bool
	IgnoreIfHeader map[string]string
	HashFunc       ETagHashFunc
}

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
		IgnoreIfHeader: map[string]string{
			"Upgrade": "websocket",
		},
		HashFunc: func(data []byte) string {
			hash := md5.Sum(data)
			return fmt.Sprintf("%x", hash)
		},
	}
}

func (t *ETagConfig) SetMinBodyLength(length int) *ETagConfig {
	t.MinBodyLength = length
	return t
}

func (t *ETagConfig) AddMethod(method string) *ETagConfig {
	t.EnableMethod[method] = true
	return t
}

func (t *ETagConfig) AddStatus(status int) *ETagConfig {
	t.EnableStatus[status] = true
	return t
}

func (t *ETagConfig) AddIgnoreHeader(header string, value string) *ETagConfig {
	t.IgnoreIfHeader[header] = value
	return t
}

func (t *ETagConfig) SetHashFunc(hashFunc ETagHashFunc) *ETagConfig {
	t.HashFunc = hashFunc
	return t
}
