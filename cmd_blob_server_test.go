//
// Simple testing of the blob-server
//
//
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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

//
// Test uploading a file.
//
func TestBlobUpload(t *testing.T) {

	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	//
	// Prepare the handler
	//
	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}", UploadHandler).Methods("POST")

	// Get the test-server
	ts := httptest.NewServer(router)
	defer ts.Close()

	//
	// The content we're going to upload
	//
	content := []byte("Content goes here, honest")

	//
	// Test a valid name
	//
	url := ts.URL + "/blob/123456"

	resp, err := http.Post(url, url, bytes.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response-body %v\n", err)
	}

	//
	// Test the upload succeeded.
	//
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	//
	// Get the result of the upload, which is a JSON string
	// which should contain:  "status":"OK"
	//
	result := fmt.Sprintf("%s", body)

	//
	// Test that it does contain that.
	//
	if !strings.Contains(result, "\"status\":\"OK\"") {
		t.Fatalf("Unexpected body: '%s'", result)
	}

	//
	// Now the file should exist
	//
	if !STORAGE.Exists("123456") {
		t.Errorf("Exists('123456') failed, post-upload!")
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)

}

//
// Test uploading a file with a bogus ID
//
func TestBlobUploadBogusID(t *testing.T) {
	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	//
	// Prepare the handler
	//
	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}", UploadHandler).Methods("POST")

	// Get the test-server
	ts := httptest.NewServer(router)
	defer ts.Close()

	//
	// The content we're going to upload
	//
	content := []byte("Content goes here, honest")

	//
	// Test a bogus name.
	//
	url := ts.URL + "/blob/foo-bar"

	resp, err := http.Post(url, url, bytes.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response-body %v\n", err)
	}

	red := fmt.Sprintf("%s", body)

	if status := resp.StatusCode; status != http.StatusInternalServerError {
		t.Errorf("Unexpected status-code: %v", status)
	}

	if red != "Alphanumeric IDs only.\n" {
		t.Errorf("Unexpected status-code: %v", red)
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)
}

//
// Test a round-trip of upload & download.
//
func TestBlobRoundTrip(t *testing.T) {

	//
	// Create a temporary directory.
	//
	p, _ := ioutil.TempDir(os.TempDir(), "prefix")

	//
	// Init the filesystem storage-class - defined in `cmd_blob_server.go`
	//
	STORAGE = new(FilesystemStorage)
	STORAGE.Setup(p)

	//
	// Prepare the handlers for upload & download
	//
	router := mux.NewRouter()
	router.HandleFunc("/blob/{id}", UploadHandler).Methods("POST")
	router.HandleFunc("/blob/{id}", GetHandler).Methods("GET")

	//
	// Get the test-server
	//
	ts := httptest.NewServer(router)
	defer ts.Close()

	//
	// The content & filename we're going to upload.
	//
	bcontent := []byte("Content goes here, honest")
	filename := "123456"

	//
	// Before the upload the file won't exist.
	//
	if STORAGE.Exists(filename) {
		t.Errorf("Exists() was true, pre-upload!")
	}

	//
	// Perform the upload.
	//
	url := ts.URL + "/blob/" + filename

	//
	// Here we're performing the HTTP-POST, but we're doing
	// it in a roundabout fashion so that we can see a "mime-type".
	//
	// This is done via the X-Mime-Type header.
	//
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bcontent))
	req.Header.Add("X-Mime-Type", "binary/steve")
	req.Header.Add("X-File-Name", filename)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//
	// If we received an error then we're done.
	//
	if err != nil {
		t.Errorf("Failed to read response-body %v\n", err)
	}

	//
	// Test the upload succeeded.
	//
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	//
	// Get the result of the upload, which is a JSON string
	// which should contain:  "status":"OK"
	//
	result := fmt.Sprintf("%s", body)
	if !strings.Contains(result, "\"status\":\"OK\"") {
		t.Fatalf("Unexpected body: '%s'", result)
	}

	//
	// Now the file should exist in the storage-directory
	//
	if !STORAGE.Exists(filename) {
		t.Errorf("Exists('') failed, post-upload!")
	}

	//
	// Now we've:
	//
	// 1. Uploaded the file.
	// 2. Tested that this succeeded.
	//
	// We need to download it
	//
	//
	req, err = http.NewRequest("GET", "/blob/"+filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status-code: %v", status)
	}

	//
	// The fake X-Mime-Type header we sent at submission-time
	// should be visible now.
	//
	if rr.Header().Get("Content-Type") != "binary/steve" {
		t.Errorf("Unexpected content-type: %v", rr.Header().Get("Content-Type"))
	}
	if rr.Header().Get("X-File-Name") != filename {
		t.Errorf("Custom filename wasn't found")
	}

	//
	// Check the response body is what we expect.
	//
	if rr.Body.String() != string(bcontent) {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), bcontent)
	}

	//
	// Cleanup the storage-point
	//
	os.RemoveAll(p)

}

//
// Test our 404-handler (!)
//
func TestMissing(t *testing.T) {
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(MissingHandler)

	//
	// Each of the methods we test.
	//
	methods := []string{"GET", "POST", "HEAD"}

	for _, method := range methods {

		//
		// Each of the paths we test
		//
		paths := []string{"/robots.txt", "/favicon.ico", "/moi.kissa"}

		for _, path := range paths {

			req, err := http.NewRequest(method, path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("Unexpected status-code for %s - %s: %v", method, path, status)
			}
		}
	}
}
