/*
 * main.go - Syncing utility.
 *
 */

package main

import (
	"../blobservers"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var verbose *bool

//
// Does the given server contain the given object?
//
func HasObject(server string, object string) bool {

	response, err := http.Head(server + "/blob/" + object)
	if err != nil {
		fmt.Printf("Error Fetching: %s/blob/%s\n", server, object)
		return false
	} else {
		if response.StatusCode == 200 {
			fmt.Printf("\tObject %s is present on %s\n", object, server)
			return true
		} else {
			fmt.Printf("\tObject %s is missing on %s\n", object, server)
			return false
		}
	}
}

//
// Mirror the given object.
//
func MirrorObject(src string, dst string, obj string) {
	// NOP
	if *verbose {
		fmt.Printf("\t\tMirroring %s from %s to %s\n",
			obj, src, dst)
	}
}

//
// Sync the members of the set we're given.
//
func SyncGroup(servers []blobservers.BlobServer) {
	//
	// If we're being verbose show the members
	//
	if *verbose {
		for _, s := range servers {
			fmt.Printf("\tGroup member: %s\n", s.Location)
		}
	}

	//
	// For each server - download the content-list here
	//
	//   key is server-name
	//   val is array of strings
	//
	objects := make(map[string][]string)

	for _, s := range servers {
		type list_strings []string
		var tmp list_strings

		file, err := ioutil.ReadFile("test.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &tmp)
		if err != nil {
			log.Fatal(err)
		}

		//
		// OK we have a string-array - we'll add that
		// to our list.
		//
		for _, i := range tmp {

			if *verbose {
				fmt.Printf("\tServer %s has %s\n", s.Location, i)
			}
			objects[s.Location] = append(objects[s.Location], i)
		}
	}

	//
	// Right we have a list of servers.
	//
	// For each server we also have the list of objects
	// that they contain.
	//
	for _, server := range servers {

		// The objects on this server
		var obs = objects[server.Location]

		// For each object.
		for _, i := range obs {

			//
			//  Mirror the object to every server that is not itself
			//
			for _, mirror := range servers {

				//
				// Ensure that src != dst.
				//
				if mirror.Location != server.Location {

					// If the object is missing.
					if !HasObject(mirror.Location, i) {
						MirrorObject(server.Location, mirror.Location, i)
					}
				}

			}
		}
	}
}

/**
 * Entry point to our code.
 */
func main() {

	//
	// Parse our command-line arguments.
	//
	blob := flag.String("blob-server", "", "Comma-separated list of blob-servers to contact.")
	verbose = flag.Bool("verbose", false, "Be verbose?")
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
	// Show the blob-servers.
	//
	if *verbose {
		fmt.Printf("\t% 10s - %s\n", "group", "server")
		for _, entry := range blobservers.Servers() {
			fmt.Printf("\t% 10s - %s\n", entry.Group, entry.Location)
		}
	}

	fmt.Println("TODO: Write a replication tool!")

	//
	// Get a list of groups.
	//
	for _, entry := range blobservers.Groups() {

		if *verbose {
			fmt.Printf("Syncing group: %s\n", entry)
		}

		//
		// For each group, get the members, and sync them.
		//
		SyncGroup(blobservers.GroupMembers(entry))
	}
}