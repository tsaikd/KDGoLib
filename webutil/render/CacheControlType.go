package render

import "github.com/tsaikd/KDGoLib/enumutil"

// CacheControlType main enum type
type CacheControlType int8

// List all valid enum
const (
	CacheControlNoStore CacheControlType = 1 + iota
	CacheControlNoCache
	CacheControlPrivate
	CacheControlPublic
)

var factoryCacheControl = enumutil.NewEnumFactory().
	Add(CacheControlNoStore, "no-store").
	Add(CacheControlNoCache, "no-cache").
	Add(CacheControlPrivate, "private").
	Add(CacheControlPublic, "public").
	Build()

func (t CacheControlType) String() string {
	return factoryCacheControl.String(t)
}
