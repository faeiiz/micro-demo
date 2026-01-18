package main

import (
	"bytes"
	"encoding/json"
	"fmt" // Added fmt for logging
	"io"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	From     string `json:"from"`
	Received string `json:"received"`
}

func main() {

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "redis:6379",
	// 	Password: "", // no password by default
	// 	DB:       0,
	// })
	// 1. Initialize Environment Variables
	service2URL := os.Getenv("SERVICE2_URL")
	if service2URL == "" {
		service2URL = "http://service2:8081/process"
	}
	fmt.Printf("1. service2URL initialized: %s\n", service2URL)

	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		// 2. Request Received
		fmt.Println("2. Received request on /proxy")

		if r.Method != http.MethodPost {
			fmt.Println("2.1. Error: Method not allowed")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req map[string]interface{}
		// 3. Decoding Request Body

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Printf("3. Error decoding JSON: %v\n", err)
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// 4. Preparing Request for Service 2
		body, _ := json.Marshal(req)
		fmt.Printf("4. Marshaled JSON for forwarding: %s\n", string(body))

		// 5. Calling Service 2
		fmt.Printf("5. Sending POST request to: %s\n", service2URL)
		resp, err := http.Post(service2URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("5.1. Error calling service2: %v\n", err)
			http.Error(w, "failed to call service2", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		fmt.Printf("5.2. Received response status from service2: %d\n", resp.StatusCode)

		// 6. Responding to Client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		written, err := io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("6. Error copying response body: %v\n", err)
		}
		fmt.Printf("6.1. Proxying complete. Bytes written: %d\n", written)
	})

	// 7. Server Start
	fmt.Println("7. service1 starting and listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
