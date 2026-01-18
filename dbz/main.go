package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Character struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Power int   `json:"power"`
}

var chars = []Character{
	{ID: "1", Name: "Goku", Power: 9001},
	{ID: "2", Name: "Vegeta", Power: 8500},
	{ID: "3", Name: "Gohan", Power: 7000},
}

func main() {
	http.HandleFunc("/chars", handleAll)
	http.HandleFunc("/chars/", handleByID)
	http.HandleFunc("/chars/search", handleSearch)
	log.Println("dbz running on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(chars)
}

func handleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "id missing", http.StatusBadRequest)
		return
	}
	id := parts[1]
	for _, c := range chars {
		if c.ID == id {
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.NotFound(w, r)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	res := []Character{}
	for _, c := range chars {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(name)) {
			res = append(res, c)
		}
	}
	json.NewEncoder(w).Encode(res)
}
