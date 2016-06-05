/*
 * main.go - Syncing utility.
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/skx/sos/blobservers"
	"strings"
)

/**
 * Entry point to our code.
 */
func main() {

	//
	// Parse our command-line arguments.
	//
	blob := flag.String("blob-server", "", "Comma-separated list of blob-servers to contact.")
	verbose := flag.Bool("verbose", false, "Be verbose?")
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
	if ( *verbose ) {
		fmt.Printf("\t% 10s - %s\n", "group", "server")
		for _, entry := range blobservers.Servers() {
			fmt.Printf("\t% 10s - %s\n", entry.Group, entry.Location)
		}
	}

	fmt.Println("TODO: Write a replication tool!")
}
