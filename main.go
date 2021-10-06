package main

import (
	"log"
	"net/http"

	"github.com/kubideh/kubesearch/hello"
)

func main() {
	log.Println("Hello World!")

	if err := http.ListenAndServe(":8080", http.HandlerFunc(hello.Handler)); err != nil {
		log.Fatalln(err)
	}
}
