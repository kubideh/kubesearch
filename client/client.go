package client

import (
	"flag"
	"fmt"

	"github.com/kubideh/kubesearch/search/api"
)

// ConfigureDefault configures and returns a new Client.
func ConfigureDefault() Client {
	flags := CreateImmutableClientFlags()
	flags.Parse()

	return Create(flags)
}

// Create returns Client objects.
func Create(flags ImmutableClientFlags) Client {
	return Client{
		flags: flags,
	}
}

// Client provides everything needed to run kubectl-search.
type Client struct {
	flags ImmutableClientFlags
}

func (c Client) serverEndpoint() string {
	return "http://" + c.flags.Server()
}

// Run creates a client that uses the given server endpoint to
// queryString for Kubernetes objects.
func (c Client) Run() error {
	result, err := api.Search(c.serverEndpoint(), queryString())

	fmt.Println(result)

	return err
}

func queryString() string {
	return flag.Arg(0)
}
