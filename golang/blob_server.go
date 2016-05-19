/*
 * blob_server.go
 *
 * Trivial rewrite of the blob-server in #golang.
 *
 * TODO:
 *   Output JSON on-upload.
 *
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//
// The location our data is stored at
//
var ROOT = "./data"

/**
 * Called via GET /blob/XXXXXX
 */
func GetHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	vars := mux.Vars(req)
	id := vars["id"]

	fname := ROOT + "/" + id

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		http.NotFound(res, req)
		return
	}

	http.ServeFile(res, req, fname)
}

/**
 * Fallback handler, returns 404 for all requests.
 */
func MissingHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "404 - content is not hosted here.")
}

/**
 * List the blob-IDs of all blobs we know about.
 */
func ListHandler(res http.ResponseWriter, req *http.Request) {

	var list []string

	files, _ := ioutil.ReadDir(ROOT)
	for _, f := range files {
		list = append(list, f.Name())
	}

	// Now output that to JSON
	mapB, _ := json.Marshal(list)
	fmt.Fprintf(res, string(mapB))
}

/**
 * Upload a file to to the public-root.
 */
func UploadHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	vars := mux.Vars(req)
	id := vars["id"]
	fname := ROOT + "/" + id

	/**
         * Get the incoming body and write to the given
         * file-name.
	*/
	var outfile *os.File
	if outfile, err = os.Create(fname); nil != err {
		status = http.StatusInternalServerError
		return
	}

	if _, err = io.Copy(outfile, req.Body); nil != err {
		status = http.StatusInternalServerError
		return
	}

	//
	// Write out the redirection - using the host
	// scheme, and the end-point of the new upload.
	//
	hostname := req.Host
	scheme := "http"

	if strings.HasPrefix(req.Proto, "HTTPS") {
		scheme = "https"
	}
	if req.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	//
	// Output the get-link.
	//
	// TODO: Write out JSON data about the upload.
	//
	res.Write([]byte(scheme + "://" + hostname + "/blob/" + id + "\n"))
}

//
// Entry-point.
//
func main() {

	host := flag.String("host", "127.0.0.1", "The IP to listen upon")
	port := flag.Int("port", 3001, "The port to bind upon")
	store := flag.String("store", "data", "The location to write the data  to")

	flag.Parse()
	ROOT = *store

	/* Create a router */
	router := mux.NewRouter()

	/* Get a previous upload */
	router.HandleFunc("/blob/{id}", GetHandler).Methods("GET")

	/* Post a new one */
	router.HandleFunc("/blob/{id}", UploadHandler).Methods("POST")

	/* List blobs we know about */
	router.HandleFunc("/blob", ListHandler).Methods("GET")

	/* Error-Handler - Return a 404 on all requests */
	router.PathPrefix("/").HandlerFunc(MissingHandler)

	/* Load the routers beneath the server root */
	http.Handle("/", router)

	/* Launch the server */
	fmt.Printf("Launching the server on http://%s:%d\n", *host, *port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic(err)
	}
}
