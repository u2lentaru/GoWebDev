package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	//urls := []string{"https://tass.ru", "https://rbc.ru", "https://ria.ru"}

	// for _, resurl := range strSearch(os.Args[1], urls) {
	// 	fmt.Println(resurl)
	// }

	pk := "https://yadi.sk/i/03bE933n3PqpG2"
	yaurl := "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=" + pk
	//yaurl := "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=https://yadi.sk/i/03bE933n3PqpG2"
	getYandexFile(yaurl)

}

func strSearch(str string, urls []string) []string {

	resultUrls := []string{}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return resultUrls
		}
		fmt.Println(url)
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return resultUrls
		}
		fmt.Println(strings.Count(string(body), str))
		if strings.Contains(string(body), str) {
			resultUrls = append(resultUrls, url)
		}
	}

	fmt.Println("Pages contains '" + str + "':")

	return resultUrls
}

func getYandexFile(yaurl string) {
	fmt.Println(yaurl)
	resp, err := http.Get(yaurl)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Body)

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
