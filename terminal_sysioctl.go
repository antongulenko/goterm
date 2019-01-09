// +build !windows,!plan9,!solaris

package goterm

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var warnOnce sync.Once

func getWinsize() (*unix.Winsize, error) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, os.NewSyscallError("GetWinsize", err)
	}

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(unix.TIOCGWINSZ),
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
