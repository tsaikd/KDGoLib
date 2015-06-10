package mmfile

import (
	"errors"
	"syscall"
	"unsafe"
)

type memfile struct {
	ptr  uintptr
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
	fmap, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READWRITE, 0, uint32(fsize), nil)
	if err != nil {
		return
	}
	defer syscall.CloseHandle(fmap)

	ptr, err := syscall.MapViewOfFile(fmap, syscall.FILE_MAP_READ|syscall.FILE_MAP_WRITE, 0, 0, uintptr(fsize))
	if err != nil {
		return nil, err
	}
	defer func() {
		if recover() != nil {
			mf = nil
			err = errors.New("Failed option a file")
		}
	}()
	mf = &memfile{
		ptr:  ptr,
		size: fsize,
		data: (*[1 << 30]byte)(unsafe.Pointer(ptr))[:fsize],
	}
	return
}

func (mf *memfile) Size() int64 {
	return mf.size
}

func (mf *memfile) Data() []byte {
	return mf.data
}

func (mf *memfile) Close() (err error) {
	return syscall.UnmapViewOfFile(mf.ptr)
}
