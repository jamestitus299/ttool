//go:build !windows

package term

import (
	"os"
	"syscall"
	"unsafe"
)

// Inject Command into the terminal's input buffer
func Inject(command string) error {
	fd := int(os.Stdin.Fd())
	for i := 0; i < len(command); i++ {
		c := command[i]
		_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSTI, uintptr(unsafe.Pointer(&c)))
		if err != 0 {
			return err
		}
	}
	return nil
}
