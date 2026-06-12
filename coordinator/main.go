package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchResult struct {
	ID    string
	Score int
}

func main(){
	resp, err := http.Get(
		"http://localhost:5001/search?q=grpc",
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	var results []SearchResult
	json.NewDecoder(resp.Body).Decode(&results)
	fmt.Println(results)
}