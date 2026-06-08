package main

import (
	"fmt"
	"strings"
)

type Document struct{
	ID		string
	Text	string
}

type SearchResult struct{
	ID		string
	Score	int
}

func rankResults(index map[string][]string, query string) []SearchResult {
	words := strings.Fields(query)

	scores := make(map[string]int)

	for _, word := range words{
		for _, id := range index[word]{
			scores[id]++
		}
	}

	results := []SearchResult{}

	for id, score := range scores {
		results = append(results, SearchResult{
			ID: id,
			Score: score,
		})
	}

	// sort here
	for i:=0; i<len(results); i++{
		for j:=i+1; j<len(results); j++{
			if results[j].Score > results[i].Score{
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results
}

func union (a []string, b []string) [] string {
	result := [] string{}

	for _, idA := range a{
		result = append(result, idA)
	}

	for _, idB := range b{
		found := false

		for _, existingId := range result {
			if existingId == idB {
				found = true
			}
		}
		if !found {
			result = append(result, idB)
		}
	}

	return result
}

func intersect(a []string, b []string) [] string{
	result := []string{}

	for _, idA := range a{
		for _, idB := range b{
			if idA == idB {
				result = append(result, idA)
			}
		}
	}

	return result
}

func searchOR(index map[string][]string, query string) [] string{
	words := strings.Fields(query)

	if len(words) == 0 {
		return []string{}
	}

	result := index[words[0]]

	for _, word := range words[1:]{
		result = union(result, index[word])
	}
	return result
}

func searchAND(index map[string][]string, query string) [] string{
	words := strings.Fields(query)

	if len(words) == 0 {
		return []string{}
	}

	result := index[words[0]]

	for _, word := range words[1:]{
		result = intersect(result, index[word])
	}
	return result
}



func main(){
	docs:= []Document{
		{ID: "1", Text: "go go go is fast"},
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



	rankedResults := rankResults(index, "go fast")

	fmt.Println(rankedResults)

	for _, result := range rankedResults {
		for _, doc := range docs {
			if doc.ID == result.ID {
				fmt.Println(doc.Text, "Score:", result.Score)
			}
		}
	}

	// orResults := searchOR(index, "go")

	// andResults := searchAND(index, "go")

	// for _, id := range orResults{
	// 	for _, doc := range docs{
	// 		if doc.ID == id {
	// 			fmt.Println(doc.Text)
	// 		}
	// 	}
	// }
}