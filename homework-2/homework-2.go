package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", firstHandle)
	router.HandleFunc("/setuser", setUsername)
	router.HandleFunc("/getuser", getUsername)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func firstHandle(wr http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		wr.Write([]byte("Hello world!"))
		return
	}
	fmt.Println("Request:" + string(body))

	type SearchStruct struct {
		Search string   `json:"search"`
		Sites  []string `json:"sites"`
	}
	searchStruct := SearchStruct{}

	err = json.Unmarshal(body, &searchStruct)
	if err != nil {
		fmt.Println(err)
		wr.Write([]byte("Hello world!"))
		return
	}

	resultUrls := make([]string, 0)

	for _, url := range searchStruct.Sites {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			//return resultUrls
		}
		fmt.Printf("There %v string counts\n", strings.Count(string(body), searchStruct.Search))
		if strings.Contains(string(body), searchStruct.Search) {
			resultUrls = append(resultUrls, url)
		}
	}

	type ResultStruct struct {
		Sites []string `json:"sites"`
	}
	resultStruct := ResultStruct{resultUrls}
	resJSON, err := json.Marshal(resultStruct)
	if err != nil {
		fmt.Println(err)
		wr.Write([]byte("Hello world!"))
		return
	}

	fmt.Printf("%s JSON\n", resJSON)
	wr.Write([]byte(resJSON))

}

func setUsername(wr http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(wr, "Hello, %s!", req.URL.Query().Get("name"))
	cun := http.Cookie{
		Name: "username",
		//Value:  "thedroppedcookiehasgoldinit",
		Value:  req.URL.Query().Get("name"),
		MaxAge: 3600}
	http.SetCookie(wr, &cun)

	wr.Write([]byte("new cookie created!\n"))
}

func getUsername(wr http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(wr, "Hello, %s!", req.URL.Query().Get("name"))
	cun, err := req.Cookie("username")
	if err != nil {
		wr.Write([]byte("error in reading cookie : " + err.Error() + "\n"))
	} else {
		value := cun.Value
		wr.Write([]byte("cookie has : " + value + "\n"))
	}
}
