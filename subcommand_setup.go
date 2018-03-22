//+build !test

package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

//
// This file contains the boiler-plate for the subcommands.
//
// It is deliberately setup into a single file, which is excluded
// from the tests via the magic-comment at the top of the file,
// so that unit-test and test-coverage statistics do not include
// this code - which cannot meaningfully be tested.
//

//
// Options which may be set via flags for the "api-server" subcommand.
//
type apiServerCmd struct {
	host    string
	blob    string
	dport   int
	uport   int
	dump    bool
	verbose bool
}

//
// Glue
//
func (*apiServerCmd) Name() string     { return "api-server" }
func (*apiServerCmd) Synopsis() string { return "Launch an API-server." }
func (*apiServerCmd) Usage() string {
	return `API-server :
  Launch an API-server to handle the upload/download of objects.
`
}

//
// Flag setup
//
func (p *apiServerCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.host, "api-host", "0.0.0.0", "The IP to listen upon.")
	f.StringVar(&p.blob, "blob-server", "", "Comma-separated list of blob-servers to contact.")
	f.IntVar(&p.dport, "download-port", 9992, "The port to bind upon for downloading objects.")
	f.IntVar(&p.uport, "upload-port", 9991, "The port to bind upon for uploading objects.")
	f.BoolVar(&p.dump, "dump", false, "Dump configuration and exit?")
	f.BoolVar(&p.verbose, "verbose", false, "Show more output from the API-server.")
}

//
// Entry-point - pass control to the API-server setup function.
//
func (p *apiServerCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	api_server(*p)
	return subcommands.ExitSuccess
}

//
// Options which may be set via flags for the "blob-server" subcommand.
//
type blobServerCmd struct {
	store string
	port  int
	host  string
}

//
// Glue
//
func (*blobServerCmd) Name() string     { return "blob-server" }
func (*blobServerCmd) Synopsis() string { return "Launch a blob-server." }
func (*blobServerCmd) Usage() string {
	return `blob-server :
  Launch a blob-server to handle the back-end storage
`
}

//
// Flag setup
//
func (p *blobServerCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.host, "host", "127.0.0.1", "The IP to listen upon")
	f.IntVar(&p.port, "port", 3001, "The port to bind upon")
	f.StringVar(&p.store, "store", "data", "The location to write the data  to")
}

//
// Entry-point.
//
func (p *blobServerCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	blob_server(*p)
	return subcommands.ExitSuccess
}

//
// Options which may be set via flags for the "replicate" subcommand.
//
type replicateCmd struct {
	blob    string
	verbose bool
}

//
// Glue
//
func (*replicateCmd) Name() string     { return "replicate" }
func (*replicateCmd) Synopsis() string { return "Trigger replication." }
func (*replicateCmd) Usage() string {
	return `replication :
  Trigger a single run of the replication/balancing operation.
`
}

//
// Flag setup
//
func (p *replicateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.blob, "blob-server", "", "Comma-separated list of blob-servers to contact.")
	f.BoolVar(&p.verbose, "verbose", false, "Be more verbose?")
}

//
// Entry-point - invoke the main replication-routine.
//
func (p *replicateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	replicate(*p)
	return subcommands.ExitSuccess
}

//
// Options which may be set via flags for the "version" subcommand.
//
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
// Entry-point.
//
func (p *versionCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	showVersion(*p)
	return subcommands.ExitSuccess
}
