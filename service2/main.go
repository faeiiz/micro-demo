package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	From     string `json:"from"`
	Received string `json:"received"`
}

func main() {
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		resp := Response{
			From:     "service2",
			Received: req.Message,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("service2 listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
