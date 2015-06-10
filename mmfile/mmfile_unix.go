// +build !windows

package mmfile

import "syscall"

type memfile struct {
	size int64
	data []byte
}

func Open(filename string) (mf MMFile, err error) {
	f, fs, err := prepareFile(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fsize := fs.Size()
	mem, err := syscall.Mmap(int(f.Fd()), 0, int(fsize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return
	}

	mf = &memfile{
		size: fsize,
		data: mem,
	}
	return
}

func (mf *memfile) Size() int64 {
	return mf.size
}

func (mf memfile) Data() []byte {
	return mf.data
}

func (mf memfile) Close() (err error) {
	return syscall.Munmap(mf.data)
}
