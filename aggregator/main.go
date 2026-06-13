package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchResult struct {
	ID    string  `json:"ID"`
	Score float64 `json:"Score"`
}

func fetchShard(url string) []SearchResult {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error calling shard:", err)
		return nil
	}
	defer resp.Body.Close()

	var results []SearchResult

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		fmt.Println("decode error:", err)
		return nil
	}

	return results
}
