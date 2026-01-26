package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID    string
	Name  string
	Email string
}

// static token -> user mapping (replace with DB later)
var tokens = map[string]User{
	"token-admin-123": {ID: "1", Name: "Admin", Email: "admin@example.com"},
	"token-mimi-456":  {ID: "2", Name: "Mimi", Email: "mimi@example.com"},
}

func main() {
	dsn := "host=khazana user=demo password=demo dbname=demo port=5432 sslmode=disable"
	fmt.Println(dsn)
	// Backend service addresses (use k8s DNS or docker-compose names)
	pehchanURL := envOr("PEHCHAN_URL", "http://pehchan:8081")
	pokemonURL := envOr("POKEMON_URL", "http://pokemon:8082")
	dbzURL := envOr("DBZ_URL", "http://dbz:8083")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		proxyTo(w, r, pehchanURL, "/auth/login")
	})

	mux.HandleFunc("/api/", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// route based on prefix after /api/
		path := strings.TrimPrefix(r.URL.Path, "/api/")
		switch {
		case strings.HasPrefix(path, "pehchan"):
			proxyTo(w, r, pehchanURL, strings.TrimPrefix(path, "pehchan"))

		case strings.HasPrefix(path, "pokemon"):
			proxyTo(w, r, pokemonURL, strings.TrimPrefix(path, "pokemon"))

		case strings.HasPrefix(path, "dbz"):
			proxyTo(w, r, dbzURL, strings.TrimPrefix(path, "dbz"))

		default:
			http.Error(w, "unknown api route", http.StatusNotFound)
		}
	}))

	addr := envOr("LISTEN_ADDR", ":8080")
	log.Printf("service1 starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// authMiddleware enforces presence of Authorization header "Bearer <token>".
// If token valid, calls the handler with the resolved user.
func authMiddleware(next func(http.ResponseWriter, *http.Request, User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		user, ok := tokens[token]
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		// Add minimal debugging header (optional)
		w.Header().Set("X-Authenticated-User", user.ID)
		next(w, r, user)
	}
}

// proxyTo forwards the incoming request to targetBase (e.g., http://pehchan:8081)
// The forwarded path is targetBase + targetPath (targetPath must begin with / or be "")
func proxyTo(w http.ResponseWriter, r *http.Request, targetBase string, targetPath string) {
	targetURL := strings.TrimRight(targetBase, "/") + targetPath
	// Build new request
	var body io.Reader
	if r.Body != nil {
		buf, _ := io.ReadAll(r.Body)
		r.Body.Close()
		body = bytes.NewReader(buf)
		// reset r.Body in case downstream needs it (not strictly necessary here)
		r.Body = io.NopCloser(bytes.NewReader(buf))
	}
	req, err := http.NewRequest(r.Method, targetURL, body)
	if err != nil {
		http.Error(w, "failed to build request", http.StatusInternalServerError)
		return
	}
	// copy headers
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	// optional: set a short timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("failed to call backend %s: %v", targetURL, err)
		http.Error(w, msg, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// copy status code and headers
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

var jwtSecret = []byte("my-secret-key")

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
