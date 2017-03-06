// +build !windows,!plan9,!solaris

package goterm

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

var warnOnce sync.Once

func getWinsize() (*winsize, error) {
	ws := new(winsize)

	var _TIOCGWINSZ int64

	switch runtime.GOOS {
	case "linux":
		_TIOCGWINSZ = 0x5413
	case "darwin":
		_TIOCGWINSZ = 1074295912
	}

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		warnOnce.Do(func() {
			fmt.Fprintln(os.Stderr, "goterm.getWinsize Error:", os.NewSyscallError("GetWinsize", errno))
		})
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}
