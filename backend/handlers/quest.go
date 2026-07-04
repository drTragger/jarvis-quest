package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const passwordHash = "bd2db407d9a55fe38c78ab85cae88528f6fdfcb441912c6b3e27afe583ee1448"

var stepAnswers = map[string]string{
	"step1": "040eb379271f55d0ddbc0abca1ab71b5307b7bd898592b4be8b3e60dc6d844e2",
	"step2": "6ed8919ce20490a5e3ad8630a4fab69475297abd07db73918dd5f36fcfaeb11b",
}

type hintRule struct {
	DelaySeconds int
	MaxAttempts  int
	Text         string
}

var stepHints = map[string]hintRule{
	"step1": {
		DelaySeconds: 120,
		MaxAttempts:  2,
		Text:         "Антонім слова \"stupid\". Одинадцять літер. Перша — I.",
	},
	"step2": {
		DelaySeconds: 300,
		MaxAttempts:  3,
		Text:         "Слухати марно. Потрібен інструмент, що малює звук частотою і часом. П'ять літер.",
	},
}

type QuestHandler struct {
	Store *ProgressStore
}

func hashMatches(input, expected string) bool {
	normalized := strings.ToLower(strings.TrimSpace(input))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:]) == expected
}

func (h *QuestHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	unlocked := h.Store.Get()
	json.NewEncoder(w).Encode(map[string]any{"unlocked": unlocked})
}

func (h *QuestHandler) StartStep(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	stepID := r.URL.Query().Get("step_id")
	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	h.Store.EnsureStepStarted(stepID)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (h *QuestHandler) ValidateAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
		StepID   string `json:"step_id"`
		Answer   string `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(req.Password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"correct": false})
		return
	}

	h.Store.EnsureStepStarted(req.StepID)

	expected, exists := stepAnswers[req.StepID]
	if !exists || !hashMatches(req.Answer, expected) {
		h.Store.IncrementAttempts(req.StepID)
		json.NewEncoder(w).Encode(map[string]bool{"correct": false})
		return
	}

	stepNum, _ := strconv.Atoi(strings.TrimPrefix(req.StepID, "step"))
	unlocked := h.Store.Advance(stepNum + 1)

	json.NewEncoder(w).Encode(map[string]any{
		"correct":  true,
		"unlocked": unlocked,
	})
}

func (h *QuestHandler) GetHint(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	stepID := r.URL.Query().Get("step_id")
	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	rule, exists := stepHints[stepID]
	if !exists {
		json.NewEncoder(w).Encode(map[string]any{"eligible": false})
		return
	}

	startedAt := h.Store.EnsureStepStarted(stepID)
	attempts := h.Store.GetAttempts(stepID)
	elapsed := time.Since(time.Unix(startedAt, 0))

	eligible := elapsed >= time.Duration(rule.DelaySeconds)*time.Second || attempts >= rule.MaxAttempts

	if !eligible {
		remaining := int((time.Duration(rule.DelaySeconds)*time.Second - elapsed).Seconds())
		json.NewEncoder(w).Encode(map[string]any{
			"eligible":          false,
			"seconds_remaining": remaining,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"eligible": true,
		"hint":     rule.Text,
	})
}
