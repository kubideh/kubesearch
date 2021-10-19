package main

import (
	"log"
	"net/http"

	_ "github.com/kubideh/kubesearch/hello"
)

func main() {
	log.Println("Hello World!")

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Fatalln(err)
	}
}
