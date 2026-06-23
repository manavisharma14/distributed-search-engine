package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     500,
		IdleConnTimeout:     90 * time.Second,
	},
}

type SearchResult struct {
	ID    string  `json:"ID"`
	Score float64 `json:"Score"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	defer func() {

		fmt.Println("global search took:", time.Since(start))

	}()

	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	resultsChan := make(chan []SearchResult, 4)

	go func() {
		resultsChan <- fetchShard(
			"http://shard1:5001/search?q=" + query,
		)
	}()

	go func() {
		resultsChan <- fetchShard(
			"http://shard2:5002/search?q=" + query,
		)
	}()

	go func() {
		resultsChan <- fetchShard(
			"http://shard3:5003/search?q=" + query,
		)
	}()

	go func() {
		resultsChan <- fetchShard(
			"http://shard4:5004/search?q=" + query,
		)
	}()

	shard1 := <-resultsChan
	shard2 := <-resultsChan
	shard3 := <-resultsChan
	shard4 := <-resultsChan

	results := []SearchResult{}

	results = append(results, shard1...)
	results = append(results, shard2...)
	results = append(results, shard3...)
	results = append(results, shard4...)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	const K = 20

	if len(results) > K {
		results = results[:K]
	}

	// for i, result := range results {
	// 	fmt.Printf(
	// 		"%d. Doc %s Score %.4f\n",
	// 		i+1,
	// 		result.ID,
	// 		result.Score,
	// 	)
	// }

	// if len(results) > 0 {
	// 	fmt.Println("global top result:", results[0])
	// }

	// fmt.Println("search took:", time.Since(start))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func fetchShard(url string) []SearchResult {

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error calling shard:", err)
		return nil
	}

	// fmt.Println(
	// 	"fetch",
	// 	url,
	// 	"took",
	// 	time.Since(start),
	// )

	defer resp.Body.Close()

	var results []SearchResult

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		fmt.Println("decode error:", err)
		return nil
	}
	return results
}

func main() {
	http.HandleFunc("/search", searchHandler)
	fmt.Println("aggregator running on :8080")
	http.ListenAndServe(":8080", nil)
}
