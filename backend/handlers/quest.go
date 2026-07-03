package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
)

var stepAnswers = map[string]string{
	"step1": "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8", // "password"
}

func ValidateAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		StepID string `json:"step_id"`
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	normalized := strings.ToLower(strings.TrimSpace(req.Answer))
	hash := sha256.Sum256([]byte(normalized))
	expected, exists := stepAnswers[req.StepID]

	w.Header().Set("Content-Type", "application/json")

	if !exists || hex.EncodeToString(hash[:]) != expected {
		json.NewEncoder(w).Encode(map[string]bool{"correct": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"correct": true})
}
