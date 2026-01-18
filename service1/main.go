package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// 1. Define the routing table (Internal K8s DNS names)
	services := map[string]string{
		"/pehchan": "http://pehchan-service:8080",
		"/pokemon": "http://pokemon-service:8080",
		"/dbz":     "http://dbz-service:8080",
	}

	// 2. Define the main handler (The entry point)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// --- STEP 1: AUTHENTICATION ---
		token := r.Header.Get("X-Auth-Token") // Or "Authorization"
		if token == "" {
			fmt.Println("Auth failed: Missing token")
			http.Error(w, "Missing Authentication Token", http.StatusUnauthorized)
			return
		}

		// Static check (We will move this to Postgres later)
		if token != "my-secret-key" {
			fmt.Println("Auth failed: Invalid token")
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		// --- STEP 2: DYNAMIC ROUTING ---
		var targetURL *url.URL
		found := false

		for prefix, serviceAddr := range services {
			if strings.HasPrefix(r.URL.Path, prefix) {
				var err error
				targetURL, err = url.Parse(serviceAddr)
				if err != nil {
					log.Printf("Target URL parse error: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Optional: Strip the prefix
				// Example: /pokemon/pikachu -> /pikachu
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Route not found for: %s\n", r.URL.Path)
			http.Error(w, "Service Not Found", http.StatusNotFound)
			return
		}

		// --- STEP 3: REVERSE PROXY ---
		// This replaces your manual http.Post logic.
		// It copies all headers, body, and method (GET/POST/PUT) automatically.
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		fmt.Printf("Forwarding request to: %s%s\n", targetURL.Host, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Gateway service1 starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
