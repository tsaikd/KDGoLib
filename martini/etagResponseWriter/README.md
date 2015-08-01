# martini etag http header response middleware

## Usage

```go
package main

import (
	"github.com/go-martini/martini"
	"github.com/tsaikd/KDGoLib/martini/etagResponseWriter"
)

func main() {
	m := martini.Classic()
	m.Use(etagResponseWriter.ETagResponseWriter(etagResponseWriter.NewETagConfig()))
	m.Get("/", func() string {
		return "Hello World!"
	})
	m.Run()
}
```

## Reference

* https://github.com/cgarvis/martini-etag
