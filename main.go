package main

import ("fmt"
		"html/template"
		"net/http"
		"strings"
	)

var index = make(map[string][]string)
var docs []Document

func helloHandle(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	tmpl.Execute(w, nil)
}

func main(){

	docs = []Document{
		{ID: "1", Text: "go go go is fast"},
		{ID: "2", Text: "go is simple"},
		{ID: "3", Text: "search engines use indexes"},
	}

	for _, doc := range docs{
		words := strings.Fields(doc.Text)

		for _, word := range words {
			index[word] = append(index[word], doc.ID)
		}
	}

	http.HandleFunc("/", helloHandle)
	
	fmt.Println("server running on :8080")

	http.ListenAndServe(":8080", nil)
}