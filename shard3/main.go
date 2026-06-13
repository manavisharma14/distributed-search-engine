package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SearchResult struct {
	ID    string
	Score int
}

type Document struct {
	ID   string
	Text string
}

var index map[string][]string
var documents []Document

func generateDocuments(startID, n int) {
	keywords := []string{
		"grpc",
		"distributed",
		"golang",
		"concurrency",
		"redis",
		"cache",
		"docker",
		"kubernetes",
		"microservices",
		"scaling",
	}

	for i := startID; i < startID+n; i++ {
		doc := Document{
			ID: strconv.Itoa(i),
			Text: keywords[i%len(keywords)] + " " +
				keywords[(i+1)%len(keywords)] + " " +
				keywords[(i+2)%len(keywords)],
		}

		documents = append(documents, doc)
	}
}

func buildIndex() {
	index = make(map[string][]string)

	for _, doc := range documents {
		words := strings.Fields(doc.Text)

		for _, word := range words {
			index[word] = append(index[word], doc.ID)
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	scores := make(map[string]int)
	query := r.URL.Query().Get("q")
	words := strings.Fields(query)

	for _, word := range words {
		ids := index[word]
		for _, id := range ids {
			scores[id]++
		}
	}
	fmt.Println(scores)

	results := []SearchResult{}

	for id, score := range scores {
		results = append(results, SearchResult{
			ID:    id,
			Score: score,
		})
	}

	json.NewEncoder(w).Encode(results)
}

func main() {
	generateDocuments(20001, 10000)
	buildIndex()
	http.HandleFunc("/search", searchHandler)
	fmt.Println("shard server running on :5003")
	http.ListenAndServe(":5003", nil)
}
