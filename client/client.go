package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubideh/kubesearch/search/api"
)

// Create configures and returns a new Client.
func Create() Client {
	flags := newFlags()
	flags.Parse()

	return New(flags)
}

// New returns Client objects.
func New(flags immutableFlags) Client {
	return Client{
		flags: flags,
	}
}

// Client provides everything needed to run kubectl-search.
type Client struct {
	flags immutableFlags
}

func (c Client) serverEndpoint() string {
	return "http://" + c.flags.Server()
}

// Run creates a client that uses the given server endpoint to
// query for Kubernetes objects.
func (c Client) Run() {
	result, err := api.Search(c.serverEndpoint(), query())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func query() string {
	return flag.Arg(0)
}
