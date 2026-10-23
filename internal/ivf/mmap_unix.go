//go:build !windows

package ivf

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
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

	data, err = unix.Mmap(int(f.Fd()), 0, int(size),
		unix.PROT_READ, unix.MAP_PRIVATE|unix.MAP_POPULATE)
	if err != nil {
		return nil, 0, fmt.Errorf("mmap: %w", err)
	}

	_ = unix.Madvise(data, unix.MADV_RANDOM)
	_ = unix.Madvise(data, unix.MADV_WILLNEED)
	_ = unix.Mlock(data)

	return data, 0, nil
}

func munmapFile(data []byte, h uintptr) error {
	if len(data) > 0 {
		_ = unix.Munlock(data)
		return unix.Munmap(data)
	}
	return nil
}
