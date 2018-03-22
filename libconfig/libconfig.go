//
// The code in this file relates to the blob-servers.
//
// The list of blob-servers is read from /etc/sos.conf + ~/.sos.conf
//
// The simple version of this file is a literal list of blob-servers:
//
//    http://node1.example.com:3333/
//    http://node2.example.com:3333/
//    http://node3.example.com:3333/
//
// When you start to scale to more than a couple of servers this becomes
// an INI-file, with one section per logical-grouping of blob-servers:
//
//  [1]
//  http://node1.example.com:1234/
//  http://mirror1-1.example.com:1234/
//  http://mirror1-2.example.com:1234/
//
//  [2]
//  http://node2.example.com:1234/
//  http://mirror2-1.example.com:1234/
//  http://mirror2-1.example.com:1234/
//
// See `SCALING.md` for the rationale behind this setup.
//

package libconfig

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

//
// This struct represents a blob-server.
//
// The blob-server has:
//
//  *  A location (host:port).
//  *  A group to which it belongs.
//
type BlobServer struct {
	Location string
	Group    string
}

//
// The list of servers we've identified.
//
var servers []BlobServer

//
// Return the list of servers we've discovered.
//
func Servers() []BlobServer {
	return (servers)
}

//
// Return the name of each group we have defined.
//
func Groups() []string {
	groups := []string{}
	for _, entry := range servers {
		found := false
		for _, a := range groups {
			if entry.Group == a {
				found = true
			}
		}
		if found == false {
			groups = append(groups, entry.Group)

		}
	}
	return (groups)
}

//
// Return the members of the given group
//
func GroupMembers(group string) []BlobServer {
	ret := []BlobServer{}

	for _, entry := range servers {
		if entry.Group == group {
			ret = append(ret, entry)
		}
	}
	return (ret)

}

//
// This returns a priority-ordered list of servers which will be used
// for uploads/downloads.
//
// Remember that we have potentially N groups.  We want to pick the first
// server from each group, then the second server from each group, etc, etc.
//
// Given input
//  group:1 host:host1
//  group:1 host:host2
//  group:2 host:host1
//  group:2 host:host2
//  group:3 host:host1
//  group:3 host:host2
//  group:3 host:host3
//
// The output should be:
//
//  group:1 host1
//  group:2 host1
//  group:3 host1
//  group:1 host2
//  group:2 host2
//  group:3 host2
//  group:3 host3
//
// This means we take three accesses to hit a server in each group, rather
// than five.  Similar savings will add up when there are more groups and
// servers.
//
func OrderedServers() []BlobServer {
	var res []BlobServer

	//
	// Create a copy of `servers`, the global list of all
	// known blob-servers
	//
	tmp := make([]BlobServer, len(servers))
	copy(tmp, servers)

	//
	// Get the names of each distinct group.
	//
	groups := make(map[string]int)
	for _, entry := range tmp {
		groups[entry.Group] += 1
	}

	//
	//  Now we'll walk over our servers, until we've processed
	// all of them.
	//
	processed := len(tmp)

	//
	//  NOTE: Dummy loop index.
	//
	for processed != 0 {

		//
		// Have we seen the given group on this loop?
		//
		seen := make(map[string]bool)

		//
		// For each server
		//
		for o, entry := range tmp {

			//
			// If this is the first time we've seen this group
			//
			if (!seen[entry.Group]) && (entry.Location != "#") {

				// Record we've seen it now.
				seen[entry.Group] = true

				// Add the appropriate entry to our result.
				res = append(res, entry)

				// Ensure we don't see it again.
				tmp[o].Location = "#"

				// We've processed a fresh host
				processed -= 1
			}
		}
	}

	// Return the magically reshuffled set of servers.
	return (res)
}

//
// Inititalize our list of servers.
//
func InitServers() {
	ServersLoad("/etc/sos.conf")
	ServersLoad(os.ExpandEnv("$HOME/.sos.conf"))
}

//
// Add an entry to our server-list.
//
func AddServer(group string, entry string) {
	tmp := BlobServer{Location: entry, Group: group}
	servers = append(servers, tmp)
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

	//
	// Here we temporarily save the servers we've found.
	//
	// We do this so that if we come across a line which
	// looks like an ini-file then we can just abort.
	//
	tmp := []string{}

	//
	// Does this look like an INI-file?
	//
	ini_file := false

	//
	// Read the input-file line by line
	//
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		//
		// Does this line just have http://.... ?
		//
		// If so add the line to the temporary array.
		//
		if strings.HasPrefix(line, "http") {
			tmp = append(tmp, line)
		}

		//
		// Otherwise this might be an INI-file.
		//
		if strings.Contains(line, "[") {
			// This is an INI-file
			ini_file = true
		}
	}

	//
	// Is this an INI-file?
	//
	if ini_file {
		//
		// Parse it as an INI-file
		//
		cfg, err := ini.Load(file)
		if err != nil {
			panic(err)
		}

		//
		//  Process each section
		//
		sections := cfg.Sections()
		for _, name := range sections {

			if name.Name() != "DEFAULT" {

				//
				//  Get the keys.
				//
				keys := cfg.Section(name.Name()).Keys()

				for _, val := range keys {

					//
					// For each entry add to the server-list.
					//
					AddServer(name.Name(), val.String())
				}
			}
		}
		return

	} else {
		//
		// This was not an INI-file, so we just create a new
		// entry for each line we read natively.
		//
		// We'll call the (anonymous) group "default".
		for _, s := range tmp {
			AddServer("default", s)
		}
	}
}

//
//  Test-code.
//
// func main() {
//
// 	InitServers()
//
// 	//
// 	// Default, post-parse order
// 	//
// 	fmt.Printf("Default order:\n")
// 	fmt.Printf("\t% 12s - %s\n", "group", "hosts" )
// 	for _, entry := range Servers() {
// 		fmt.Printf("\t% 12s - %s\n", entry.group, entry.location)
// 	}
//
// 	//
// 	// Now the the order we prefer
// 	//
// 	fmt.Printf("Improved order:\n")
// 	fmt.Printf("\t% 12s - %s\n", "group", "hosts" )
// 	for _, entry := range OrderedServers() {
// 		fmt.Printf("\t% 12s - %s\n", entry.group, entry.location)
// 	}
//
// }
//
