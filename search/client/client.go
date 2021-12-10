package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubideh/kubesearch/search/api"
)

func Run() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Use \"kubectl search [flags] <query>\".")
	}
	server := flag.String("server", "localhost:8080", "the address and port of the Kubesearch server")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "You must specify a query.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	result, err := api.Search("http://"+*server, flag.Arg(0))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}
