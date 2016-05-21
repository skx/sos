//
// The code in this file relates to the blob-servers.
//
// The list of blob-servers is read from /etc/sos.conf + ~/.sos.conf
//
// The general form of the line is:
//
//    http://$host:$port/ [weight]
//
// For example:
//
//    http://example.com:3333/
//    http://example.com:3333/ 10
//    http://example.com:3333/ 30
//
//  If no weight is specified then the default will be used "0".
//

package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

//
// This struct represents a blob-server.
//
// The blob-server has:
//
//  *  A location (host:port).
//
//  *  A weight.
//
// Weighting is implemented by ordering the hosts by the weight
// field, such that the bigger numbers come first.
//
type BlobServer struct {
	location string
	weight   int
}

//
// Implement a sorting interface, of `weight`.
//
type ByWeight []BlobServer

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].weight > a[j].weight }

//
// The list of servers we've identified.
//
var servers []BlobServer

//
// Return the list of servers we've read, sorted by weight order.
//
func Servers() []BlobServer {
	return (servers)
}

//
// Inititalize our list of servers.
//
func InitServers() {
	ServersLoad("/etc/sos.conf")
	ServersLoad(os.ExpandEnv("$HOME/.sos.conf"))

	// Now run the sorting
	sort.Sort(ByWeight(servers))
}

//
// Add an entry to our server-list manually.
//
func AddServer(entry string) {

	if strings.HasPrefix(entry, "http") {

		// Split the line into fields,
		// to see if there is a trailing number.
		parts := strings.Fields(entry)

		//
		// If there is then the number is the weight.
		//
		weight := 0
		if len(parts) > 1 {
			weight, _ = strconv.Atoi(parts[1])
		}

		// Add the entry.
		if len(parts) > 0 {
			tmp := BlobServer{location: parts[0], weight: weight}
			servers = append(servers, tmp)
		}

		// Now ensure our list is ordered by weight.
		sort.Sort(ByWeight(servers))
	}
}

//
// Read/Parse the list of servers from the specified file.
//
func ServersLoad(file string) {
	inFile, err := os.Open(file)
	if err != nil {
		return
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		//
		// We've read a line.
		//
		line := scanner.Text()

		AddServer(line)
	}
}
