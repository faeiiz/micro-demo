package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	http.HandleFunc("/auth/login", handleLogin)
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

var jwtSecret = []byte("my-secret-key")

func handleLogin(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	json.NewDecoder(r.Body).Decode(&body)

	if body.Email != "mimi@example.com" || body.Password != "123" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":   "2",
		"name":  "Mimi",
		"email": "mimi@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": ss,
	})
}
