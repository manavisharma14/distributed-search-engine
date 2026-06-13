package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchResult struct {
	ID    string  `json:"ID"`
	Score float64 `json:"Score"`
}
