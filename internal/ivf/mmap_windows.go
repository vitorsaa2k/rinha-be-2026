//go:build windows

package ivf

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func mmapFile(path string) (data []byte, h uintptr, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	size := fi.Size()
	if size == 0 {
		return nil, 0, fmt.Errorf("empty file: %s", path)
	}

	mh, err := windows.CreateFileMapping(
		windows.Handle(f.Fd()),
		nil,
		windows.PAGE_READONLY,
		0,
		uint32(size),
		nil,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("CreateFileMapping: %w", err)
	}

	ptr, err := windows.MapViewOfFile(
		mh,
		windows.FILE_MAP_READ,
		0, 0,
		uintptr(size),
	)
	if err != nil {
		windows.CloseHandle(mh)
		return nil, 0, fmt.Errorf("MapViewOfFile: %w", err)
	}

	data = unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(size))
	return data, uintptr(mh), nil
}

func munmapFile(data []byte, h uintptr) error {
	if len(data) > 0 {
		if err := windows.UnmapViewOfFile(uintptr(unsafe.Pointer(&data[0]))); err != nil {
			return err
		}
	}
	return windows.CloseHandle(windows.Handle(h))
}
