package apimgr

// Definition contain API detail
type Definition struct {
	Version     uint8
	Name        string
	FullName    string
	Description string
	Method      string
	Pattern     string
	Handlers    []interface{}
	Request     interface{}
	Extra       interface{}
}
