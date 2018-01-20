//
// Entry-point to the puppet-summary service.
//

package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

//
// Setup our sub-commands and use them.
//
func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	//subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&apiServerCmd{}, "")
	subcommands.Register(&blobServerCmd{}, "")
	subcommands.Register(&versionCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
