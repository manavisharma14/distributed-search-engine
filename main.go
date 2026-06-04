package main

import (
	"fmt"
	"strings"
)

type Document struct{
	ID		string
	Text	string
}

func search(index map[string][]string, query string) []string{
	words := strings.Fields(query)
	for _, word := range words {
		fmt.Println(word, "->", index[word])
	}
	return nil
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

	results := search(index, "go fast")
	fmt.Println(results)

	for _, id := range results{
		for _, doc := range docs{
			if doc.ID == id{
				fmt.Println(doc.Text)
			}
		}

	}

	// for word, docIDs := range index {
	// 	fmt.Printf("%s -> %v\n", word, docIDs)
	// }



	// fmt.Println("Documents: ")

	// for _, doc := range docs{
	// 	fmt.Printf("ID: %s | Text: %s\n", doc.ID, doc.Text)
	// }
}