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
	"io/ioutil"
	"os"
	"strings"
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

	SOSChroot(connection)
}

//
// Get the contents of a given ID.
//
func (this *FilesystemStorage) Get(id string) (*[]byte, map[string]string) {

	//
	// If the file is missing we return nil.
	//
	if _, err := os.Stat(id); os.IsNotExist(err) {
		return nil, nil
	}

	//
	// Read the file contents
	//
	x, err := ioutil.ReadFile(id)

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
	meta_data, err := ioutil.ReadFile(id + ".json")

	//
	// If we did then we can decode it
	//
	if err == nil {
		json.Unmarshal([]byte(meta_data), &meta)
		return &x, meta
	} else {

		//
		// There was a failure to unmarshal the meta-data
		// just return the actual data.
		//
		return &x, nil
	}
}

//
// Store the specified data against the given file.
//
func (this *FilesystemStorage) Store(id string, data []byte, params map[string]string) bool {

	//
	// Write out the data.
	//
	err := ioutil.WriteFile(id, data, 0644)

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
		err = ioutil.WriteFile(id+".json", encoded, 0644)

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
		name := f.Name()

		if !strings.HasSuffix(name, ".json") {
			list = append(list, name)
		}
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
