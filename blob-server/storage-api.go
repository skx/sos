//
// Storage-abstraction.
//
//  This file contains an abstract interface for storing/retrieving
// data, as well as a concrete implementation for doing that over
// a local filesystem.
//
//  It is possible that other users would be interested in storing
// data inside MySQL, Postgres, Redis, or similar.  To do that
// should involve only implementing the `StorageHandler` interface
// and changing `main.go` to construct the new instance.
//
//

package main

import (
	"io/ioutil"
	"os"
	"syscall"
)

//
//  This is the virtual interface for a storage class.
//
type StorageHandler interface {

	//
	// Setup method.
	//
	// If you're implement a database-based system the string here
	// might be used to specify host/user/password, etc.
	//
	Setup(connection string)

	//
	// Retrieve the contents of a blob by ID.
	//
	Get(id string) *[]byte

	//
	// Store some data against the given ID.
	//
	Store(id string, data []byte) bool

	//
	// Get all known IDs.
	//
	Existing() []string

	//
	// Does the given ID exist?
	//
	Exists(id string) bool
}

//
// A concrete type which implements our interface.
//
type FilesystemStorage struct {
}

//
// Setup method to ensure we have a data-directory.
//
func (this *FilesystemStorage) Setup(connection string) {

	//
	// If the data-directory does not exist create it.
	//
	os.MkdirAll(connection, 0755)

	//
	// Now try to secure ourselves
	//
	syscall.Chdir(connection)
	syscall.Chroot(connection)
}

//
// Get the contents of a given ID.
//
func (this *FilesystemStorage) Get(id string) *[]byte {

	if _, err := os.Stat(id); os.IsNotExist(err) {
		return nil
	}

	// Read the file contents
	x, err := ioutil.ReadFile(id)

	// If there was no error return the contents.
	if err != nil {
		return nil
	} else {
		return &x
	}
}

//
// Store the specified data against the given file.
//
func (this *FilesystemStorage) Store(id string, data []byte) bool {

	err := ioutil.WriteFile(id, data, 0644)
	if err != nil {
		return false
	} else {
		return true
	}
}

//
// Get all known IDs.
//
// We assume we've been chdir() + chroot() into the data-directory
// so we just need to read the filenames we can find.
//
func (this *FilesystemStorage) Existing() []string {
	var list []string

	files, _ := ioutil.ReadDir(".")
	for _, f := range files {
		list = append(list, f.Name())
	}
	return list
}

//
// Does the given ID exist (as a file)?
//
func (this *FilesystemStorage) Exists(id string) bool {
	if _, err := os.Stat(id); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
