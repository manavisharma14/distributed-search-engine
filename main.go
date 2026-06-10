// single-node search engine 
package main

import (
		"bufio"
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

type DisplayResult struct {
	Text		string
	Score		int	
}
 
type PageData struct {
	Query 		string
	Results 	[]DisplayResult
}

func buildDisplayResults(query string) [] DisplayResult{
	results := rankResults(index, query)
	displayResults := []DisplayResult{}

	for _, result := range results {
		for _, doc := range docs {
			if doc.ID == result.ID {
				displayResults = append(displayResults, DisplayResult{
					Text: doc.Text,
					Score: result.Score,
				})
			}
		}
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
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	documents := []Document{}

	scanner := bufio.NewScanner(file)


	for scanner.Scan() {
		text := scanner.Text()

		doc := Document{

		Text: text,
	}
	documents = append(documents, doc)

	}

	return documents
}

func main(){

	files, err := os.ReadDir("documents")

	if err != nil {
		fmt.Println(err)
		return
	}


	docs = []Document{}

	for _, file := range files {
		go loadDocuments("documents/" + file.Name())
	}

	fmt.Println("building index now")

	id := 1
	for i := range docs {
		docs[i].ID = strconv.Itoa(id)
		id++
	}

	for _, doc := range docs{
		words := strings.Fields(doc.Text)

		for _, word := range words {
			index[word] = append(index[word], doc.ID)
		}
	}

	http.HandleFunc("/", helloHandle)
	http.HandleFunc("/search", apiSearchHandle)
	
	

	go fmt.Println("hello from goroutine!")

	fmt.Println("server running on :8080")

	http.ListenAndServe(":8080", nil)
}