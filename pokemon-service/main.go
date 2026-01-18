package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Pokemon struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

var pokemons = []Pokemon{
	{ID: "1", Name: "Pikachu", Type: "Electric"},
	{ID: "2", Name: "Charmander", Type: "Fire"},
	{ID: "3", Name: "Bulbasaur", Type: "Grass"},
}

func main() {
	http.HandleFunc("/pokemons", handleAll)
	http.HandleFunc("/pokemons/", handleByID)
	http.HandleFunc("/pokemons/search", handleSearch)
	log.Println("pokemon running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(pokemons)
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
	for _, p := range pokemons {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
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
	res := []Pokemon{}
	for _, p := range pokemons {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			res = append(res, p)
		}
	}
	json.NewEncoder(w).Encode(res)
}
