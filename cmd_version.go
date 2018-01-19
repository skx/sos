//
// Show our version - This uses a level of indirection for our test-case
//

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
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

type versionCmd struct {
	verbose bool
}

//
// Glue
//
func (*versionCmd) Name() string     { return "version" }
func (*versionCmd) Synopsis() string { return "Show our version." }
func (*versionCmd) Usage() string {
	return `version :
  Report upon our version, and exit.
`
}

//
// Flag setup
//
func (p *versionCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.verbose, "verbose", false, "Show go version the binary was generated with.")
}

//
// Show the version - using the "out"-writer.
//
func showVersion(verbose bool) {
	fmt.Fprintf(out, "%s\n", version)
	if verbose {
		fmt.Fprintf(out, "Built with %s\n", runtime.Version())
	}
}

//
// Entry-point.
//
func (p *versionCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	showVersion(p.verbose)
	return subcommands.ExitSuccess
}
