package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = map[string]User{
	"1": {ID: "1", Name: "Admin User"},
	"2": {ID: "2", Name: "Guest User"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" || id == "all" {
			json.NewEncoder(w).Encode(users)
			return
		}
		if user, ok := users[id]; ok {
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "User Not Found", 404)
		}
	})
	http.ListenAndServe(":8080", nil)
}
