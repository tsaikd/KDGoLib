# martini enhanced return handler

## Usage

```go
package main

import (
	"fmt"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/tsaikd/KDGoLib/martini/errorJson"
)

func main() {
	m := martini.Classic()
	m.Map(errorJson.ReturnErrorProvider())
	m.Use(render.Renderer())
	m.Get("/", func() err {
		return fmt.Errorf("error")
	})
	m.Run()
}
```
