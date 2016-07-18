package futil

import (
	"os"

	"github.com/tsaikd/KDGoLib/errutil"
)

// IsExist return true if path exist
//
// https://gist.github.com/mattes/d13e273314c3b3ade33f
// if _, err := os.Stat("/path/to/whatever"); err == nil {
// 	// path/to/whatever exists
// }
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsNotExist return true if path not exist
//
// https://gist.github.com/mattes/d13e273314c3b3ade33f
// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
// 	// path/to/whatever does not exist
// }
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
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
