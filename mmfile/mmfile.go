package mmfile

import (
	"os"
)

type MMFile interface {
	Size() int64
	Data() []byte
	Close() (err error)
}

func prepareFile(filename string) (file *os.File, fi os.FileInfo, err error) {
	if file, err = os.OpenFile(filename, os.O_RDWR, 0644); err != nil {
		return
	}

	if fi, err = file.Stat(); err != nil {
		return
	}

	return
}
