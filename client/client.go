package client

import (
	"flag"
	"fmt"

	"github.com/kubideh/kubesearch/search/api"
)

// Create configures and returns a new Client.
func Create() Client {
	flags := NewFlags()
	flags.Parse()

	return New(flags)
}

// New returns Client objects.
func New(flags ImmutableFlags) Client {
	return Client{
		flags: flags,
	}
}

// Client provides everything needed to run kubectl-search.
type Client struct {
	flags ImmutableFlags
}

func (c Client) serverEndpoint() string {
	return "http://" + c.flags.Server()
}

// Run creates a client that uses the given server endpoint to
// query for Kubernetes objects.
func (c Client) Run() error {
	result, err := api.Search(c.serverEndpoint(), query())

	fmt.Println(result)

	return err
}

func query() string {
	return flag.Arg(0)
}
