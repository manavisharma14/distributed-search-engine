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

func apiSearchHandle(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")
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

	json.NewEncoder(w).Encode(displayResults)
}

func helloHandle(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	
	query := r.URL.Query().Get("q")
	results := rankResults(index, query)
	displayResults := []DisplayResult{}

	for _, result := range results{
		for _, doc := range docs{
			if doc.ID == result.ID {
				displayResults = append(displayResults, DisplayResult{
					Text: doc.Text,
					Score: result.Score,
				})
			}
		}
	}

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

	id := 1
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)

		doc := Document{
		ID: strconv.Itoa(id),
		Text: text,
	}
	documents = append(documents, doc)
	id++
	}

	return documents
}

func main(){

	files, err := os.ReadDir("documents")

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		fmt.Println("documents/" + file.Name())

	}

	docs = []Document{}

	for _, file := range files {
		docs = append(docs, loadDocuments("documents/" + file.Name())...)
	}

	for _, doc := range docs{
		words := strings.Fields(doc.Text)

		for _, word := range words {
			index[word] = append(index[word], doc.ID)
		}
	}

	http.HandleFunc("/", helloHandle)
	http.HandleFunc("/search", apiSearchHandle)
	
	fmt.Println("server running on :8080")

	http.ListenAndServe(":8080", nil)
}