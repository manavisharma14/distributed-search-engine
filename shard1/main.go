package main

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
)

type Document struct{
	ID		string
	Text	string
}

type SearchResult struct{
	ID		string
	Score	int
}

var index = map[string][]string{
	"grpc": 		{"1"},
	"distributed":	{"1"},
	"golang": 		{"2"},
	"concurrency": 	{"2"},
}

func searchHandler(w http.ResponseWriter, r *http.Request){
	scores := make(map[string]int)
	query := r.URL.Query().Get("q")
	words := strings.Fields(query)

	for _, word := range words{
		ids := index[word]
		for _, id := range ids {
			scores[id]++
		}
	}

	results := []SearchResult{}

	for id, score := range scores {
		results = append(results, SearchResult{
			ID: id,
			Score: score,
		})
	}

	fmt.Println(results)

	fmt.Println(scores)

	json.NewEncoder(w).Encode(results)
}

func main(){
	http.HandleFunc("/search", searchHandler)
	
	fmt.Println("shard server running :5001")

	http.ListenAndServe(":5001", nil)
}