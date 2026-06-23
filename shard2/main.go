package main

import (
	"container/heap"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Posts struct {
	Rows []Post `xml:"row"`
}

type Post struct {
	ID    string `xml:"Id,attr"`
	Title string `xml:"Title,attr"`
	Body  string `xml:"Body,attr"`
}

type Document struct {
	ID   string
	Text string
}

type SearchResult struct {
	ID    string
	Score float64
}

type MinHeap []SearchResult

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].Score < h[j].Score
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(SearchResult))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

var documents []Document

var index map[string]map[string]int

func loadPostsXML(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var posts Posts

	decoder := xml.NewDecoder(file)

	err = decoder.Decode(&posts)
	if err != nil {
		return err
	}

	for _, post := range posts.Rows {

		text := strings.ToLower(
			post.Title + " " + post.Body,
		)

		documents = append(documents, Document{
			ID:   post.ID,
			Text: text,
		})
	}

	return nil
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
	start := time.Now()

	defer func() {

		fmt.Println("search took:", time.Since(start))

	}()

	matchedTerms := make(map[string]int)
	scores := make(map[string]float64)
	query := r.URL.Query().Get("q")
	words := strings.Fields(query)

	if len(words) == 0 {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

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

	const K = 20

	h := &MinHeap{}
	heap.Init(h)

	for id, score := range scores {
		if matchedTerms[id] != len(words) {
			continue
		}

		result := SearchResult{
			ID:    id,
			Score: score,
		}

		if h.Len() < K {
			heap.Push(h, result)
		} else if result.Score > (*h)[0].Score {
			heap.Pop(h)
			heap.Push(h, result)
		}
	}

	results := []SearchResult{}

	for h.Len() > 0 {
		results = append(results, heap.Pop(h).(SearchResult))
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	fmt.Printf("returned %d results\n", len(results))

	if len(results) > 0 {
		fmt.Println("top result:", results[0])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	err := loadPostsXML(
		"./data/Posts.xml",
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("documents:", len(documents))

	buildIndex()
	http.HandleFunc("/search", searchHandler)
	fmt.Println("shard server running on :5002")
	http.ListenAndServe(":5002", nil)
}
