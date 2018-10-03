package futil

import (
	"os"
	"path/filepath"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrProjectRootNotFound1 = errutil.NewFactory("project root not found from: %q")
)

// SearchProjectRoot return project root dir if possible,
// searchFromDir will set to current working directory if empty
func SearchProjectRoot(
	searchFromDir string,
	mustExist []string,
	devExist []string,
	prodExist []string,
) (dir string, err error) {
	if searchFromDir == "" {
		if searchFromDir, err = os.Getwd(); err != nil {
			return
		}
	}
	if dir, err = filepath.Abs(searchFromDir); err != nil {
		return
	}

	for {
		switch dir {
		case "", ".", "/":
			return "", ErrProjectRootNotFound1.New(nil, searchFromDir)
		}

		// must exist dir
		if allPathExist(dir, mustExist) {
			// dev
			if allPathExist(dir, devExist) {
				return
			}

			// production
			if allPathExist(dir, prodExist) {
				return
			}
		}

		dir = filepath.Dir(dir)
	}
}

func allPathExist(base string, names []string) bool {
	for _, name := range names {
		path := filepath.Join(base, name)
		if IsNotExist(path) {
			return false
		}
	}
	return true
}
