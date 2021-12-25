package client

import (
	"flag"
	"fmt"
	"os"
)

// NewFlags returns the ImmutableFlags for Client. A list of the
// flags and their defaults are now given.
//
// -server (default: localhost:8080)
func NewFlags() ImmutableFlags {
	return NewFlagsWithServer("localhost:8080")
}

// NewFlagsWithServer returns the ImmutableFlags for Client, and
// it uses the given server as a default value for the flag
// `-server`.
func NewFlagsWithServer(server string) ImmutableFlags {
	flag.Usage = usage

	return ImmutableFlags{
		server: flag.String("server", server, "the address and port of the KubeSearch server"),
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Use \"kubectl search [flags] <query>\".")
}

// ImmutableFlags is a collection of flags used to configure the
// Client. Each flag will be populated with values from the command-
// line after calling Parse().
type ImmutableFlags struct {
	server *string // server is an address and port that can be used by `http.Get`
}

// Server returns an address and port that can be used by
// `http.Get`, and it's populated by a value from the
// command-line.
func (f ImmutableFlags) Server() string {
	return *f.server
}

// Parse populates this collection of ImmutableFlags with values from the
// command-line, and it validates command-line arguments.
func (f ImmutableFlags) Parse() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		printUsageAndExitWithFailure()
	}
}

func printUsageAndExitWithFailure() {
	fmt.Fprintln(os.Stderr, "You must specify a query.")
	fmt.Fprintln(os.Stderr, "")
	flag.Usage()
	os.Exit(1)
}
