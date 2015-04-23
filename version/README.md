KDGoLib version package
=======================

## Usage

* import package

```go
import "github.com/tsaikd/KDGoLib/version"
```

* initialize version

```go
func init() {
	version.VERSION = "0.0.1"
}
```

* show version in somewhere

```go
func main() {
	fmt.Println(version.String())
}
```

* config LDFLAGS when compile time

```sh
githash="$(git rev-parse HEAD | cut -c1-8)"
buildtime="$(date +%Y-%m-%d)"

LDFLAGS="${LDFLAGS} -X github.com/tsaikd/KDGoLib/version.BUILDTIME ${buildtime}"
LDFLAGS="${LDFLAGS} -X github.com/tsaikd/KDGoLib/version.GITCOMMIT ${githash}"

go build -ldflags "${LDFLAGS}"
```
