package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubideh/kubesearch/search/api"
)

func Run() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Use \"kubectl search <query>\".")
	}

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "You must specify a query.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	result, err := api.Search("http://localhost:8080", flag.Args()[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}
