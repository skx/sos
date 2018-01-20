//
// Simple testing of the blob-server
//
//
package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

//
// Upload IDs must be alphanumeric.  Submit some bogus requests to
// ensure they fail with a suitable error-message.
//
func TestGetIDNames(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}/", GetHandler).Methods("GET")
	router.HandleFunc("/blob/{id}", GetHandler).Methods("GET")

	//
	// Table driven test - each of these should fail as not
	// matching "a-z0-9".
	//
	ids := []string{"/blob/xXx", "/blob/34l'", "/blob/a-b-c", "/blob/<fdf>"}

	for _, id := range ids {

		req, err := http.NewRequest("GET", id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Unexpected status-code: %v", status)
		}

		// Check the response body is what we expect.
		expected := "Alphanumeric IDs only.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got '%v' want '%v'",
				rr.Body.String(), expected)
		}
	}
}

//
// Test that our health end-point returns an alive-response.
//
func TestHealth(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/alive/", HealthHandler).Methods("GET")
	router.HandleFunc("/alive", HealthHandler).Methods("GET")

	ids := []string{"/alive", "/alive/"}

	for _, id := range ids {

		req, err := http.NewRequest("GET", id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status-code: %v", status)
		}

		// Check the response body is what we expect.
		expected := "alive"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got '%v' want '%v'",
				rr.Body.String(), expected)
		}
	}
}

//
// Test HEAD access to a file before/after it is created.
//
func TestHeadAccess(t *testing.T) {

	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}/", GetHandler).Methods("HEAD")
	router.HandleFunc("/blob/{id}", GetHandler).Methods("HEAD")

	//
	// Table driven test - each of these should fail as the matching
	// file is not present.
	//
	ids := []string{"/blob/foo", "/blob/foo/"}

	for _, id := range ids {

		req, err := http.NewRequest("HEAD", id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Unexpected status-code: %v", status)
		}
	}

	//
	// Now create the file `foo`.
	//
	path := filepath.Join(p, "foo")
	content := []byte("Content")
	ioutil.WriteFile(path, content, 0644)

	//
	// And repeat the accesses - we should have success now.
	//
	for _, id := range ids {

		req, err := http.NewRequest("HEAD", id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status-code: %v", status)
		}
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)
}

//
// Test that retrieving a missing blob fails appropriately.
//
func TestMissingBlob(t *testing.T) {

	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}/", GetHandler).Methods("GET")
	router.HandleFunc("/blob/{id}", GetHandler).Methods("GET")

	//
	// Table driven test - each of these should fail as the matching
	// file is not present.
	//
	ids := []string{"/blob/919aac866fb1fb107616a5e3824efc91aacb3be1", "/blob/8b55aac644e9e6f2701805584cc391ff81d3ecec"}

	for _, id := range ids {

		req, err := http.NewRequest("GET", id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Unexpected status-code: %v", status)
		}

		// Check the response body is what we expect.
		expected := "404 page not found\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got '%v' want '%v'",
				rr.Body.String(), expected)
		}
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)
}

//
// Test the blob-list
//
func TestBlobList(t *testing.T) {

	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	router := mux.NewRouter()
	router.HandleFunc("/blobs", ListHandler).Methods("GET")

	//
	// Nothing uploaded so we should get "[]"
	//
	req, err := http.NewRequest("GET", "/blobs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	// Check the response body is what we expect.
	expected := "[]"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	//
	// Now create a single file, so we can test that the listing
	// returns valid contents.
	//
	path := filepath.Join(p, "steve")
	content := []byte("Content")
	ioutil.WriteFile(path, content, 0644)

	//
	// At this point we should get a single result.
	//
	req, err = http.NewRequest("GET", "/blobs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	// Check the response body is what we expect.
	expected = "[\"steve\"]"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)
}
