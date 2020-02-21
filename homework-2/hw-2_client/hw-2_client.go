package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	type SearchStruct struct {
		Search string   `json:"search"`
		Sites  []string `json:"sites"`
	}
	searchStruct := SearchStruct{
		Search: "нефть",
		Sites: []string{
			"https://tass.ru",
			"https://rbc.ru",
			"https://ria.ru",
		},
	}

	sendJSON, err := json.Marshal(searchStruct)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s JSON\n", sendJSON)

	resp, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(sendJSON))
	if err != nil {
		fmt.Printf("%s !error!\n", err)
		//return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
