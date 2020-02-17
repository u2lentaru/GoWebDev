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
	strSearch(os.Args[1], urls)
	// fmt.Println("Pages contains '" + os.Args[1] + "':")

	// for _, resurl := range strSearch(os.Args[1], urls) {
	// 	fmt.Println(resurl)
	// }

}

func strSearch(str string, urls []string) {
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
		fmt.Println(strings.Count(string(body), str))
		if strings.Contains(string(body), str) {
			resultUrls = append(resultUrls, url)
		}
	}

	fmt.Println("Pages contains '" + str + "':")

	for _, resurl := range resultUrls {
		fmt.Println(resurl)
	}
	//return
}

// PS C:\Golang_work\src\GoWebDev> go run C:\Golang_work\src\GoWebDev\homework-1\homework-1.go нефть
// https://tass.ru
// 200 OK
// 0
// https://rbc.ru
// 200 OK
// 1
// https://ria.ru
// 200 OK
// 2
// Pages contains 'нефть':
// https://rbc.ru
// https://ria.ru

//https://yadi.sk/i/03bE933n3PqpG2
