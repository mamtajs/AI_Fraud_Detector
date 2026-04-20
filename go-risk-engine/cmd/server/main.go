package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"fraud-detection-platform/go-risk-engine/internal/model"
	"fraud-detection-platform/go-risk-engine/internal/scoring"
)

func main() {
	weightsPath := getenv("MODEL_WEIGHTS_PATH", "model/weights.json")
	scorer, err := scoring.NewScorer(weightsPath)
	if err != nil {
		log.Fatalf("unable to initialize scorer: %v", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/score", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req model.ScoreRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload"})
			return
		}

		resp := scorer.Score(req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	port := getenv("PORT", "8081")
	log.Printf("go-risk-engine listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
