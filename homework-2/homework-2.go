package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", firstHandle)
	router.HandleFunc("/user", helloUsername)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func firstHandle(wr http.ResponseWriter, req *http.Request) {
	wr.Write([]byte("Hello world!"))
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Request:" + string(body))
}

func helloUsername(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wr, "Hello, %s!", req.URL.Query().Get("name"))
}
