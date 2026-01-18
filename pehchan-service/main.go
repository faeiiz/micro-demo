package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: "1", Name: "Admin", Age: 30},
	{ID: "2", Name: "Mimi", Age: 26},
	{ID: "3", Name: "Alex", Age: 28},
}

func main() {
	http.HandleFunc("/users", handleUsers)         // GET /users
	http.HandleFunc("/users/", handleUserByID)     // GET /users/{id}
	http.HandleFunc("/users/search", handleSearch) // GET /users/search?name=
	log.Println("pehchan running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// path: /users/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "id missing", http.StatusBadRequest)
		return
	}
	id := parts[1]
	for _, u := range users {
		if u.ID == id {
			json.NewEncoder(w).Encode(u)
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
	res := []User{}
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			res = append(res, u)
		}
	}
	json.NewEncoder(w).Encode(res)
}
