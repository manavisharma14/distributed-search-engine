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

func searchAllShards(query string) []SearchResult{
	shards := []string{
		"http://localhost:5001",
		"http://localhost:5002",
	}

	resultChannel := make(chan []SearchResult)

	for _, shard := range shards{
		go func(url string){
			searchUrl := url + "/search?q=" + query
			resp, err := http.Get(searchUrl)

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
	return allResults
}

func main(){

	results := searchAllShards("grpc")
	fmt.Println(results)
	
}