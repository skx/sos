//
// Show our version - This uses a level of indirection for our test-case
//

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

//
// modified during testing
//
var out io.Writer = os.Stdout

var (
	version = "unreleased"
)

//
// Show the version - using the "out"-writer.
//
func showVersion(options versionCmd) {
	fmt.Fprintf(out, "%s\n", version)
	if options.verbose {
		fmt.Fprintf(out, "Built with %s\n", runtime.Version())
	}
}
