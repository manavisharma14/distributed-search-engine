package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type SearchResult struct {
	ID    string
	Score float64
}

type Document struct {
	ID   string
	Text string
}

var index map[string]map[string]int
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
	index = make(map[string]map[string]int)

	for _, doc := range documents {
		words := strings.Fields(doc.Text)

		for _, word := range words {
			if index[word] == nil {
				index[word] = make(map[string]int)
			}
			index[word][doc.ID]++
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	scores := make(map[string]float64)
	query := r.URL.Query().Get("q")
	words := strings.Fields(query)
	matchedTerms := make(map[string]int)

	for _, word := range words {
		docsContainingWord := index[word]

		if len(docsContainingWord) == 0 {
			continue
		}

		N := len(documents)
		df := len(docsContainingWord)

		idf := math.Log(float64(N) / float64(df))

		for docID, count := range docsContainingWord {
			tf := float64(count)
			scores[docID] += tf * idf
			matchedTerms[docID]++
		}
	}

	results := []SearchResult{}

	for id, score := range scores {
		if matchedTerms[id] != len(words) {
			continue
		}

		results = append(results, SearchResult{
			ID:    id,
			Score: score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	fmt.Printf("returned %d results\n", len(results))

	if len(results) > 0 {
		fmt.Println("top result:", results[0])
	}

	const K = 20

	if len(results) > K {

		results = results[:K]

	}

	json.NewEncoder(w).Encode(results)
}

func main() {
	generateDocuments(100001, 50000)
	buildIndex()
	http.HandleFunc("/search", searchHandler)
	fmt.Println("shard server running on :5003")
	http.ListenAndServe(":5003", nil)
}
