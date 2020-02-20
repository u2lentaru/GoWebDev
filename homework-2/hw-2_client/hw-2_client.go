package main

import (
	"encoding/json"
	"fmt"
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
	fmt.Printf("%s", sendJSON)

	//resp, err := http.Post("localhost:8080","application/json", bytes.NewBuffer(reqBody))

}
