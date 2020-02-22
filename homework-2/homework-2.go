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
	if req.URL.Path == "/favicon.ico" {
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		wr.Write([]byte(err.Error()))
		wr.WriteHeader(400)
		return
	}

	type TSearchStruct struct {
		Search string   `json:"search"`
		Sites  []string `json:"sites"`
	}
	SearchStruct := TSearchStruct{}

	err = json.Unmarshal(body, &SearchStruct)
	if err != nil {
		fmt.Println(err)
		wr.Write([]byte(err.Error()))
		wr.WriteHeader(400)
		return
	}

	type TResultStruct struct {
		Sites []string `json:"sites"`
	}
	ResultStruct := TResultStruct{search(SearchStruct.Sites, SearchStruct.Search)}
	resJSON, err := json.Marshal(ResultStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	wr.Write([]byte(resJSON))
}

func search(sites []string, str string) []string {
	resultUrls := make([]string, 0, len(sites))

	for _, url := range sites {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("There %v string counts\n", strings.Count(string(body), SearchStruct.Search))
		if strings.Contains(string(body), str) {
			resultUrls = append(resultUrls, url)
		}
	}
	return resultUrls
}

func setUsername(wr http.ResponseWriter, req *http.Request) {
	cun := http.Cookie{
		Name:   "username",
		Value:  req.URL.Query().Get("name"),
		MaxAge: 3600}
	http.SetCookie(wr, &cun)
	wr.Write([]byte("new cookie created!\n"))
}

func getUsername(wr http.ResponseWriter, req *http.Request) {
	cun, err := req.Cookie("username")
	if err != nil {
		wr.Write([]byte("error in reading cookie : " + err.Error() + "\n"))
		return
	}
	wr.Write([]byte("cookie has : " + cun.Value + "\n"))
}
