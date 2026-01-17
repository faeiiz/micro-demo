package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	From     string `json:"from"`
	Received string `json:"received"`
}

func main() {
	service2URL := os.Getenv("SERVICE2_URL")
	if service2URL == "" {
		service2URL = "http://service2:8081/process"
	}

	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		body, _ := json.Marshal(req)

		resp, err := http.Post(service2URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "failed to call service2", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w,resp.Body)
	})

	log.Println("service1 listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
