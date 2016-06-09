package futil

import (
	"os"

	"github.com/tsaikd/KDGoLib/errutil"
)

// IsExist return true if path exist
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// IsDir return true if path exist and is directory
func IsDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func() {
		errutil.Trace(f.Close())
	}()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}
