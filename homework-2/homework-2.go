package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", firstHandle)
	router.HandleFunc("/user", helloUsername)

	log.Fatal(http.ListenAndServe("8080", router))

}

func firstHandle(wr http.ResponseWriter, req *http.Request) {
	wr.Write([]byte("Hello world!"))
}

func helloUsername(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wr, "Hello %s!", req.URL.Query().Get("name"))
}
