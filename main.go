// single-node search engine 
package main

import (
		"os"
		"encoding/json"
		"fmt"
		"html/template"
		"net/http"
		"strings"
		"strconv"
		"sync"
	)

var index = make(map[string][]string)
var docs []Document
var shards []Shard

var docMap = make(map[string]Document)

type Shard struct {
	Docs	[]Document
	Index	map[string][]string
}

type DisplayResult struct {
	Text		string
	Score		int	
}
 
type PageData struct {
	Query 		string
	Results 	[]DisplayResult
}

func searchAllShards(query string) []SearchResult{
	allResults := []SearchResult{}

	for _, shard := range shards {
		shardResults := rankResult(shard.Index, query)
		allResults := appen(allResults, shardResult...)
	}
}

func buildDisplayResults(query string) [] DisplayResult{
	results := rankResults(index, query)
	displayResults := []DisplayResult{}

	for _, result := range results {
		doc := docMap[result.ID]

		displayResults = append(displayResults, DisplayResult{
			Text: doc.Text,
			Score: result.Score,
		})
	}
	return displayResults
}

func apiSearchHandle(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")


	displayResults := buildDisplayResults(query)

	json.NewEncoder(w).Encode(displayResults)
}

func helloHandle(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	
	query := r.URL.Query().Get("q")
	
	displayResults := buildDisplayResults(query)
	data := PageData{
		Query: query,
		Results: displayResults,
	}
	
	tmpl.Execute(w, data)
}

func loadDocuments(filename string) []Document{
	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	doc := Document{
		Text: string(data),
	}

	return []Document{doc}

	
}

func main(){

	files, err := os.ReadDir("documents")

	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	docChannel := make(chan []Document)
	docs = []Document{}

	for _, file := range files {
		wg.Add(1)

		go func(filename string){
			defer wg.Done()
			loadedDocs := loadDocuments(filename)
			docChannel <- loadedDocs
		}("documents/" + file.Name())
	}

	for range files {
		loadedDocs := <-docChannel
		docs = append(docs, loadedDocs...)
	}

	wg.Wait()

	fmt.Println("all workers finished")
	fmt.Println("documents loaded: ", len(docs))

	id := 1
	for i := range docs {
		docs[i].ID = strconv.Itoa(id)
		docMap[docs[i].ID] = docs[i]
		id++
	}
	
	shards = []Shard{
		{
			Docs: []Document{},
			Index: make(map[string][]string),
		},
		{
			Docs: []Document{},
			Index: make(map[string][]string),
		},
	}

	for i, doc := range docs {
		shardId := i%len(shards)

		shards[shardId].Docs = append(shards[shardId].Docs, doc)
	}

	for i, s := range shards {
		fmt.Println(
			"Shard",
			i,
			"documenrs:",
			len(s.Docs),
		)
	}

	

	for i:= range shards {
		for _, doc := range shards[i].Docs {
			words := strings.Fields(doc.Text)
			
			for _, word := range words {
				shards[i].Index[word] = append(shards[i].Index[word], doc.ID)
			}
		}
	}

	for i, shard := range shards {
		fmt.Println("shard", i)

		for word, ids := range shard.Index {
			fmt.Println(word, "->", ids)
		}
	}

	http.HandleFunc("/", helloHandle)
	http.HandleFunc("/search", apiSearchHandle)
	
	fmt.Println("server running on :8080")

	http.ListenAndServe(":8080", nil)
}