package main

import (
	"fmt"
	"strings"
)

type Document struct{
	ID		string
	Text	string
}

func search(index map[string][]string, query string) [] string{
	words := strings.Fields(query)

	if len(words) == 0 {
		return []string{}
	}

	fmt.Println(words)
	fmt.Println(words[0])
	fmt.Println(index[words[0]])

	return index[words[0]]
}

func main(){
	docs:= []Document{
		{ID: "1", Text: "go is fast"},
		{ID: "2", Text: "go is simple"},
		{ID: "3", Text: "search engines use indexes"},
	}

	index := make(map[string][]string)

	for _, doc := range docs{
		words := strings.Fields(doc.Text)

		for _, word := range words {
			index[word] = append(index[word], doc.ID)
		}
	}

	results := search(index, "go")
	fmt.Println(results)



}