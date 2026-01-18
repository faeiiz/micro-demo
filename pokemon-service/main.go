package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Pokemon struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var pokedex = map[string]Pokemon{
	"pikachu":   {Name: "Pikachu", Type: "Electric"},
	"charizard": {Name: "Charizard", Type: "Fire/Flying"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/"))
		if name == "" {
			json.NewEncoder(w).Encode(pokedex)
			return
		}
		if p, ok := pokedex[name]; ok {
			json.NewEncoder(w).Encode(p)
		} else {
			http.Error(w, "Pokemon Not Found", 404)
		}
	})
	http.ListenAndServe(":8080", nil)
}
