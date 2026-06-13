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

type Document struct {
	ID   string
	Text string
}

type SearchResult struct {
	ID    string
	Score float64
}

var documents []Document

var index map[string]map[string]int

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
		if i == startID {
			doc.Text = "grpc grpc grpc grpc distributed"
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

	// for _, word := range words {
	// 	ids := index[word]
	// 	for _, id := range ids {
	// 		scores[id]++
	// 	}
	// }

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
		}
	}

	results := []SearchResult{}

	for id, score := range scores {
		results = append(results, SearchResult{
			ID:    id,
			Score: score,
		})
	}

	fmt.Println(results)

	fmt.Println(scores)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	json.NewEncoder(w).Encode(results)
}

func main() {
	generateDocuments(1, 50000)
	fmt.Println(documents[0])
	buildIndex()
	fmt.Println("grpc count in doc1:", index["grpc"]["1"])

	fmt.Println("documents:", len(documents))

	// fmt.Println(index["grpc"])

	http.HandleFunc("/search", searchHandler)

	fmt.Println("shard server running :5001")

	http.ListenAndServe(":5001", nil)
}
