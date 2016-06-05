/*
 * main.go - API Server
 *
 * Trivial rewrite of the api-server in #golang.
 *
 * Presents two end-points, on two different ports:
 *
 *    http://127.0.0.1:9991/upload
 *
 *    http://127.0.0.1:9992/fetch/:id
 */

package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"github.com/skx/sos/blobservers"
)

/**
 * This is a helper for allowing us to consume a HTTP-body more than once.
 */
type myReader struct {
	*bytes.Buffer
}

/**
 * So that it implements the io.ReadCloser interface
 */
func (m myReader) Close() error { return nil }

//
// Upload-handler.
//
// This should attempt to upload against the blob-servers and return
// when that is complete.  If there is a failure then it should
// repeat the process until all known servers are exhausted.
//
// The retry logic is described in the file `SCALING.md` in the
// repository, but in brief there are two cases:
//
//  * All the servers are in the group `default`.
//
//  * There are N defined groups.
//
// Both cases are handled by the call to OrderedServers() which
// returns the known blob-servers in a suitable order to minimize
// lookups.  See `SCALING.md` for more details.
//
//
func UploadHandler(res http.ResponseWriter, req *http.Request) {

	//
	// We create a new buffer to hold the request-body.
	//
	buf, _ := ioutil.ReadAll(req.Body)

	//
	// Create a copy of the buffer, so that we can consume
	// it initially to hash the data.
	//
	rdr1 := myReader{bytes.NewBuffer(buf)}

	//
	// Get the SHA1 hash of the uploaded data.
	//
	hasher := sha1.New()
	b, _ := ioutil.ReadAll(rdr1)
	hasher.Write([]byte(b))
	hash := hasher.Sum(nil)

	//
	// Now we're going to attempt to re-POST the uploaded
	// content to one of our blob-servers.
	//
	// We try each blob-server in turn, and if/when we receive
	// a successful result we'll return it to the caller.
	//
	for _, s := range blobservers.OrderedServers() {

		//
		// Replace the request body with the (second) copy we made.
		//
		rdr2 := myReader{bytes.NewBuffer(buf)}
		req.Body = rdr2

		//
		// This is where we'll POST to.
		//
		url := fmt.Sprintf("%s%s%x", s.Location, "/blob/", hash)

		//
		// Build up a new request.
		//
		child, _ := http.NewRequest("POST", url, req.Body)

		//
		// Propogate any incoming X-headers
		//
		for header, value := range req.Header {
			if strings.HasPrefix(header, "X-") {
				child.Header.Set(header, value[0])
			}
		}

		//
		// Send the request.
		//
		client := &http.Client{}
		r, err := client.Do(child)

		//
		// If there was no error we're good.
		//
		if err == nil {

			//
			// We read the reply we received from the
			// blob-server and return it to the caller.
			//
			response, _ := ioutil.ReadAll(r.Body)

			if response != nil {
				fmt.Fprintf(res, string(response))
				return
			}
		}
	}

	//
	// If we reach here we've attempted our upload on every
	// known blob-server and none accepted it.
	//
	// Let the caller know.
	//
	fmt.Fprintf(res, "Upload FAILED")
	return

}

//
// Download-handler.
//
// This should attempt to download against the blob-servers and return
// when that is complete.  If there is a failure then it should
// repeat the process until all known servers are exhausted..
//
// The retry logic is described in the file `SCALING.md` in the
// repository, but in brief there are two cases:
//
//  * All the servers are in the group `default`.
//
//  * There are N defined groups.
//
// Both cases are handled by the call to OrderedServers() which
// returns the known blob-servers in a suitable order to minimize
// lookups.  See `SCALING.md` for more details.
//
//
func DownloadHandler(res http.ResponseWriter, req *http.Request) {

	//
	// The ID of the file we're to retrieve.
	//
	vars := mux.Vars(req)
	id := vars["id"]

	//
	// Strip any extension which might be present on the ID.
	//
	extension := filepath.Ext(id)
	id = id[0 : len(id)-len(extension)]

	//
	// We try each blob-server in turn, and if/when we receive
	// a successfully result we'll return it to the caller.
	//
	for _, s := range blobservers.OrderedServers() {

		//
		// Build up a request, which is a HTTP-GET
		//
		response, err := http.Get(fmt.Sprintf("%s%s%s", s.Location, "/blob/", id))
		//
		// If there was no error we're good.
		//
		if err != nil {
			fmt.Printf("Error fetching %s from %s%s%s\n",
				id, s.Location, "/blob/", id)
		} else {

			//
			// We read the reply we received from the
			// blob-server and return it to the caller.
			//
			body, _ := ioutil.ReadAll(response.Body)

			if body != nil {

				//
				// Copy any X-Header which was present
				// into the reply too.
				//
				for header, value := range response.Header {
					if strings.HasPrefix(header, "X-") {
						res.Header().Set(header, value[0])
					}
				}

				//
				// Now send back the body.
				//
				fmt.Fprintf(res, string(body))
				return
			}
		}
	}

	//
	// If we reach here we've attempted our download on every
	// known blob-server and none succeeded.
	//
	// Let the caller know.
	//
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "Object not found.")
}

/**
 * Entry point to our code.
 */
func main() {

	//
	// Parse our command-line arguments.
	//
	host := flag.String("host", "0.0.0.0", "The IP to listen upon.")
	blob := flag.String("blob-server", "", "Comma-separated list of blob-servers to contact.")
	dport := flag.Int("download-port", 9992, "The port to bind upon for downloading objects.")
	uport := flag.Int("upload-port", 9991, "The port to bind upon for uploading objects.")
	dump := flag.Bool("dump", false, "Dump configuration and exit?.")
	flag.Parse()

	//
	// If we received blob-servers on the command-line use them too.
	//
	// NOTE: blob-servers added on the command-line are placed in the
	// "default" group.
	//
	if (blob != nil) && (*blob != "") {
		servers := strings.Split(*blob, ",")
		for _, entry := range servers {
			blobservers.AddServer("default", entry)
		}
	} else {

		//
		//  Initialize the servers from our config file(s).
		//
		blobservers.InitServers()
	}

	//
	// If we're merely dumping the servers then do so now.
	//
	if *dump {
		fmt.Printf("\t% 10s - %s\n", "group", "server")
		for _, entry := range blobservers.Servers() {
			fmt.Printf("\t% 10s - %s\n", entry.Group, entry.Location)
		}
		return
	}

	//
	// Otherwise show a banner, then launch the server-threads.
	//
	fmt.Printf("Launching API-server\n")
	fmt.Printf("\nUpload service\nhttp://%s:%d/upload\n", *host, *uport)
	fmt.Printf("\nDownload service\nhttp://%s:%d/fetch/:id\n", *host, *dport)

	//
	// Show the blob-servers, and their weights
	//
	fmt.Printf("\nBlob-servers:\n")
	fmt.Printf("\t% 10s - %s\n", "group", "server")
	for _, entry := range blobservers.Servers() {
		fmt.Printf("\t% 10s - %s\n", entry.Group, entry.Location)
	}
	fmt.Printf("\n")

	//
	// Create a route for uploading.
	//
	up_router := mux.NewRouter()
	up_router.HandleFunc("/upload", UploadHandler).Methods("POST")

	//
	// Create a route for downloading.
	//
	down_router := mux.NewRouter()
	down_router.HandleFunc("/fetch/{id}", DownloadHandler).Methods("GET")

	//
	// The following code is a hack to allow us to run two distinct
	// HTTP-servers on different ports.
	//
	// It's almost sexy the way it worked the first time :)
	//
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *uport),
			up_router)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *dport),
			down_router)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
