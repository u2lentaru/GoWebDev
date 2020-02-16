package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://golang.org")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	io.Copy(os.Stdout, resp.Body)
}
