package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	urls := []string{"https://tass.ru", "https://rbc.ru", "https://ria.ru"}

	for _, resurl := range strSearch(os.Args[1], urls) {
		fmt.Println(resurl)
	}

	pubkey := "https://yadi.sk/i/03bE933n3PqpG2"
	getYandexFile("https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=" + pubkey)

}

func strSearch(str string, urls []string) []string {

	resultUrls := []string{}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return resultUrls
		}
		defer resp.Body.Close()
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
	type tyad struct {
		Href      string
		Method    string
		Templated bool
	}

	//!!!!!
	//Lower case don't work!
	// href      string
	// method    string
	// templated bool

	var yad tyad

	resp, err := http.Get(yaurl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &yad)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err = http.Get(yad.Href)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create("yad.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("File yad.pdf successfully copied!")

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
// File yad.pdf successfully copied!
