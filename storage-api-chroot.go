// +build !windows

package main

import (
	"syscall"
)

func SOSChroot(directory string) {
	syscall.Chroot(directory)
}
