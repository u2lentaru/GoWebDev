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

//func strSearch(str string, urls []string) []string
func strSearch(str string, urls []string) []string {
	//resultUrls := []string{}
	resultUrls := make([]string, 0)

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			//return resultUrls
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return resultUrls
		}
		fmt.Printf("There %v string counts\n", strings.Count(string(body), str))
		if strings.Contains(string(body), str) {
			resultUrls = append(resultUrls, url)
		}
	}
	if len(resultUrls) > 0 {
		fmt.Printf("Pages contains %v:\n", str)
	}
	return resultUrls
}

//func getYandexFile(yaurl string)
func getYandexFile(yaurl string) {
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

	type Tyad struct {
		Href      string
		Method    string
		Templated bool
	}
	//var yad tyad
	Yad := Tyad{}
	//!!!!!
	//Lower case don't work!
	// href      string
	// method    string
	// templated bool

	err = json.Unmarshal(body, &Yad)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err = http.Get(Yad.Href)
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
