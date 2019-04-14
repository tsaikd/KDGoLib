package pkgutil

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/tsaikd/KDGoLib/futil"
)

// Module go mod info
type Module struct {
	Path     string       // module path
	Version  string       // module version
	Versions []string     // available module versions (with -versions)
	Replace  *Module      // replaced by this module
	Time     *time.Time   // time version was created
	Update   *Module      // available update, if any (with -u)
	Main     bool         // is this the main module?
	Indirect bool         // is this module only an indirect dependency of main module?
	Dir      string       // directory holding files for this module, if any
	GoMod    string       // path to go.mod file for this module, if any
	Error    *ModuleError // error loading module
}

// ModuleError go mod error
type ModuleError struct {
	Err string // the error itself
}

// IsGoModDir return true if dir contains go.mod file
func IsGoModDir(dir string) bool {
	return futil.IsExist(filepath.Join(dir, "go.mod"))
}

var reJSONObj = regexp.MustCompile("{[^}]*}")

// ParseGoMod return go module and dep info in dir
func ParseGoMod(dir string) (result []Module, err error) {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = dir
	goModData, err := cmd.Output()
	if err != nil {
		return
	}

	modsData := reJSONObj.FindAll(goModData, -1)

	result = make([]Module, len(modsData))
	for i, d := range modsData {
		mod := Module{}
		if err = json.Unmarshal(d, &mod); err != nil {
			return []Module{}, err
		}
		result[i] = mod
	}

	return
}
