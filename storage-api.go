//
// Storage-abstraction.
//
// This file contains an abstract interface for storing/retrieving
// data, as well as a concrete implementation for doing that over
// a local filesystem.
//
// It is possible that other users would be interested in storing
// data inside MySQL, Postgres, Redis, or similar.  To do that
// should involve only implementing the `StorageHandler` interface
// and changing the setup to construct a new instance.
//
// We allow "data" to be read or written, by ID.
//
// We also allow (optional) meta-data to be written/retrieved alongside
// the data.  The latter is saved as a JSON file, alongside the data.
//

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// StorageHandler is the interface for a storage class.
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
	// The fetch will also return any (optional)
	// key=value parameters which were stored when
	// the content was uploaded.
	//
	Get(id string) (*[]byte, map[string]string)

	//
	// Store some data against the given ID.
	//
	// If any optional `key=value` parameters have been
	// sent then store them too, alongside the data.
	//
	Store(id string, data []byte, params map[string]string) bool

	//
	// Get all known IDs.
	//
	Existing() []string

	//
	// Does the given ID exist?
	//
	Exists(id string) bool
}

// FilesystemStorage is a concrete type which implements
// the StorageHandler interface.
type FilesystemStorage struct {
	// cwd records whether we're chrooted or not
	cwd bool

	// prefix holds our prefix directory if we didn't chroot
	prefix string
}

//
// Setup method to ensure we have a data-directory.
//
func (fss *FilesystemStorage) Setup(connection string) {

	//
	// If the data-directory does not exist create it.
	//
	os.MkdirAll(connection, 0755)

	//
	// We default to changing to the directory and chrooting
	//
	// But we can't do that when testing.
	//
	if flag.Lookup("test.v") != nil {
		fss.cwd = false
		fss.prefix = connection
		return
	}

	//
	// Now try to secure ourselves
	//
	syscall.Chdir(connection)

	//
	// If we're running our test-cases we don't chroot.
	//
	SOSChroot(connection)

	//
	// Since we're not testing all accesses will be based
	// upon the current working directory.
	//
	fss.cwd = true

}

//
// Get the contents of a given ID.
//
func (fss *FilesystemStorage) Get(id string) (*[]byte, map[string]string) {

	//
	// If we're not using the cwd we need to build up the complete
	// path to the file.
	//
	target := id
	if fss.cwd == false {
		target = filepath.Join(fss.prefix, id)
	}

	//
	// If the file is missing we return nil.
	//
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return nil, nil
	}

	//
	// Read the file contents
	//
	x, err := ioutil.ReadFile(target)

	// If there was an error return nil too.
	if err != nil {
		return nil, nil
	}

	//
	// Now attempt to return the meta-data too.
	//
	meta := make(map[string]string)

	//
	// Attempt to read the meta-data file.
	//
	metaData, err := ioutil.ReadFile(target + ".json")

	//
	// If we did then we can decode it
	//
	if err == nil {
		json.Unmarshal([]byte(metaData), &meta)
		return &x, meta
	}

	//
	// There was a failure to unmarshal the meta-data
	// just return the actual data.
	//
	return &x, nil
}

//
// Store the specified data against the given file.
//
func (fss *FilesystemStorage) Store(id string, data []byte, params map[string]string) bool {

	//
	// If we're not using the cwd we need to build up the complete
	// path to the file.
	//
	target := id
	if fss.cwd == false {
		target = filepath.Join(fss.prefix, id)
	}

	//
	// Write out the data.
	//
	err := ioutil.WriteFile(target, data, 0644)

	//
	// If there was an error we abort.
	//
	if err != nil {
		return false
	}

	//
	// If we received some optional parameters then write them
	// out too.
	//
	if len(params) != 0 {

		// Marshal to JSON.
		encoded, err := json.Marshal(params)

		//
		// If there was an error marshalling the meta-data
		// then our upload failed, even though the data
		// was stored already..
		//
		if err != nil {
			return false
		}

		// Write out to a .json-suffixed file.
		err = ioutil.WriteFile(target+".json", encoded, 0644)

		// If the data was saved but the meta-data wasn't
		// this is still a failure.
		if err != nil {
			return false
		}
	}

	//
	// Data written.
	//
	// Meta-data written, if supplied.
	//
	return true
}

// Existing returns all known IDs.
//
// We assume we've been chdir() + chroot() into the data-directory
// so we just need to read the filenames we can find.
//
func (fss *FilesystemStorage) Existing() []string {
	var list []string

	//
	// If we're not using the cwd we need to use our prefix, explicitly
	//
	target := "."
	if fss.cwd == false {
		target = fss.prefix
	}

	files, _ := ioutil.ReadDir(target)
	for _, f := range files {
		name := f.Name()

		if !strings.HasSuffix(name, ".json") {
			list = append(list, name)
		}
	}
	return list
}

// Exists tests whether the given ID exists (as a file).
func (fss *FilesystemStorage) Exists(id string) bool {

	//
	// If we're not using the cwd we need to build up the complete
	// path to the file.
	//
	target := id
	if fss.cwd == false {
		target = filepath.Join(fss.prefix, id)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}
