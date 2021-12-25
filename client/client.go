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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Use \"kubectl search [flags] <query>\".")
	}

	return Client{
		flags: flags,
	}
}

// Client provides everything needed to run kubectl-search.
type Client struct {
	flags immutableFlags
}

func (c Client) ServerEndpoint() string {
	return "http://" + c.flags.Server()
}

func (c Client) Query() string {
	return flag.Arg(0)
}

// Run creates a client that uses the given server endpoint to
// query for Kubernetes objects.
func (c Client) Run() {
	result, err := api.Search(c.ServerEndpoint(), c.Query())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}
