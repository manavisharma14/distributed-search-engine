package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type SearchResult struct {
	ID    string
	Score int
}

func searchHandler(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")
	results := searchAllShards(query)

	json.NewEncoder(w).Encode(results)

}

func searchAllShards(query string) []SearchResult{
	shards := []string{
		"http://localhost:5001",
		"http://localhost:5002",
	}

	resultChannel := make(chan []SearchResult)

	for _, shard := range shards{
		go func(url string){
			searchUrl := url + "/search?q=" + query
			client := http.Client{
				Timeout: 2*time.Second,
			}
			resp, err := client.Get(searchUrl)

			if err != nil{
				fmt.Println(err)
				resultChannel <- []SearchResult{}
				return
			}

			defer resp.Body.Close()

			var results []SearchResult
			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil{
				fmt.Println(err)
				resultChannel <- []SearchResult{}
				return
			}
			
			resultChannel <- results
		}(shard)	
	}

	allResults := []SearchResult{}
			
	for range shards {
		shardResults := <- resultChannel
		allResults = append(allResults, shardResults...)
	}
	sort.Slice(allResults, func(i,j int) bool{
		return allResults[i].Score > allResults[j].Score
	})
	return allResults
}

func main(){
	
	http.HandleFunc("/search", searchHandler)
	fmt.Println("coordinator running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil{
		fmt.Println(err)
	}
}