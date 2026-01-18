package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings" // Need this to lowercase the input
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Request struct {
	City string `json:"city"` // This tag fixes the decoding
}

// 1. Map cities to coordinates (Simple Demo DB)
var cityCoords = map[string]string{
	"berlin": "latitude=52.52&longitude=13.41",
	"london": "latitude=51.51&longitude=-0.13",
	"paris":  "latitude=48.85&longitude=2.35",
	"tokyo":  "latitude=35.68&longitude=139.69",
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		json.NewDecoder(r.Body).Decode(&req)

		// 2. Normalize input (handle "London", "LONDON", "london")
		city := strings.ToLower(req.City)
		if city == "" {
			city = "berlin"
		}

		fmt.Printf("S2-1. Processing weather for: %s\n", city)

		// 3. Check Redis
		cachedData, err := rdb.Get(ctx, city).Result()
		if err == nil {
			fmt.Println("S2-2. Redis HIT! Returning cached data.")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedData))
			return
		}

		// 4. Resolve Coordinates
		coords, ok := cityCoords[city]
		if !ok {
			// Fallback if city not in our map
			fmt.Println("City not found in map, defaulting to Berlin")
			coords = cityCoords["berlin"]
		}

		// 5. Construct Dynamic URL
		fmt.Println("S2-2. Redis MISS. Calling Open-Meteo API...")
		apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s&current_weather=true", coords)

		apiResp, err := http.Get(apiURL)
		if err != nil {
			http.Error(w, "Weather API error", http.StatusBadGateway)
			return
		}
		defer apiResp.Body.Close()

		var weatherData interface{}
		json.NewDecoder(apiResp.Body).Decode(&weatherData)
		jsonData, _ := json.Marshal(weatherData)

		rdb.Set(ctx, city, jsonData, 60*time.Second)
		fmt.Println("S2-3. Data cached in Redis.")

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	fmt.Println("S2-0. Service 2 (Weather Engine) listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
