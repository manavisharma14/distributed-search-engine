package main

import (
	"fmt"
	"net/http"

)

func main(){
	fmt.Println("shard server running on :5002")
	http.ListenAndServe(":5002", nil)
}