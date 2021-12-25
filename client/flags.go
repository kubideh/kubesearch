package client

import (
	"flag"
	"fmt"
	"os"
)

// newFlags returns the immutableFlags for Client.
func newFlags() immutableFlags {
	return immutableFlags{
		server: flag.String("server", "localhost:8080", "the address and port of the KubeSearch server"),
	}
}

// immutableFlags is a collection of flags used to configure the
// Client. Each flag will be populated with values from the command-
// line after calling Parse().
type immutableFlags struct {
	server *string // server is an address and port that can be used by `http.Get`
}

// Server returns an address and port that can be used by
// `http.Get`, and it's populated by a value from the
// command-line.
func (f immutableFlags) Server() string {
	return *f.server
}

// Parse populates this collection of immutableFlags with values from the
// command-line, and it validates command-line arguments.
func (f immutableFlags) Parse() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "You must specify a query.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}
}
