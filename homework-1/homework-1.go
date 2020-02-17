package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	urls := []string{"https://tass.ru", "https://rbc.ru", "https://ria.ru"}
	resultUrls := []string{}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(url)
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(strings.Count(string(body), os.Args[1]))
		if strings.Contains(string(body), os.Args[1]) {
			resultUrls = append(resultUrls, url)
		}
	}

	for _, resurl := range resultUrls {
		fmt.Println(resurl)
	}
}

//PS C:\Golang_work\src\GoWebDev> go run C:\Golang_work\src\GoWebDev\homework-1\homework-1.go нефть
//https://tass.ru
//200 OK
//0
//https://rbc.ru
//200 OK
//1
//https://ria.ru
//200 OK
//2
//https://rbc.ru
//https://ria.ru
