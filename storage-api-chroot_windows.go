// +build windows

package main

// SOSChroot attempts to call `chroot` with the given directory.
//
// This Windows-specific implementation is a nop.
func SOSChroot(directory string) {
	// NOP
}
