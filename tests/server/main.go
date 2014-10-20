package main

import (
	"github.com/hypermusk/hypermusk/tests/apihttpimpl"
	"log"
	"net/http"
)

func main() {

	apihttpimpl.AddToMux("/api", http.DefaultServeMux)
	log.Println("Serving on 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
