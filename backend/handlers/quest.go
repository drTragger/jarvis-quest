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
	"step3": "2c0c2cef2d45475adbef682180c2c37848f7631800e185ad16694edabcf57a48",
	"step4": "18d63be10ad544a04a22c944dee01d6d864ec69b797a58edae92e6a44ad8fdbf",
	"step5": "de194e1890197480e508de206d8c80816559d65be70f51b033ad71e6176499d9",
	"step6": "3ca862d51945e7a7961be9074c0621b471104c8fc3e415482dcb125d19c82b54",
	"step7": "1db02e2ec8f6942146e8c3f592770be279de0fa28dc98c32289c1226c4ba2d30",
}

const hintDelay = 24 * time.Hour

type hintRule struct {
	Text string
}

var stepHints = map[string]hintRule{
	"step1": {
		Text: "Антонім слова \"stupid\". Одинадцять літер. Перша — i.",
	},
	"step2": {
		Text: "Слухати марно. Потрібен інструмент, що малює звук частотою і часом. П'ять літер.",
	},
	"step3": {
		Text: "Кожна пара цифр — код символу в шістнадцятковій системі (ASCII). Переведи кожен байт у літеру.",
	},
	"step4": {
		Text: "Це слово описує саме те, чим ця перевірка підозрює, що ти є. П'ять літер, з малої.",
	},
	"step5": {
		Text: "Відкрий зображення в будь-якому фоторедакторі (навіть у стандартному додатку «Фото» на телефоні) і сильно підніми яскравість або контраст. У найтемнішому кутку зображення проявиться текст. Прочитане слово додатково зашифроване шифром Цезаря зі зсувом 13 (ROT13) — застосуй той самий зсув ще раз.",
	},
	"step6": {
		Text: "Короткий спалах — крапка, довгий — тире, довша пауза розділяє літери. Розшифруй ім'я — і згадай, як тебе називає твій друг. Вісім літер, з малої, без пробілів (Владислав ...).",
	},
	"step7": {
		Text: "Прийди до кав'ярні SVOЇ, там на тебе чекатиме команда, яку треба виконати в терміналі.",
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
	json.NewEncoder(w).Encode(map[string]any{
		"unlocked":        unlocked,
		"onboarding_seen": h.Store.GetOnboardingSeen(),
		"rules_seen":      h.Store.GetRulesSeen(),
	})
}

func (h *QuestHandler) MarkRulesSeen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(req.Password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	h.Store.MarkRulesSeen()
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (h *QuestHandler) MarkOnboardingSeen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !hashMatches(req.Password, passwordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	h.Store.MarkOnboardingSeen()
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
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

	stepNum, _ := strconv.Atoi(strings.TrimPrefix(req.StepID, "step"))
	expected, exists := stepAnswers[req.StepID]
	if !exists || stepNum != h.Store.Get() || !hashMatches(req.Answer, expected) {
		h.Store.IncrementAttempts(req.StepID)
		json.NewEncoder(w).Encode(map[string]bool{"correct": false})
		return
	}

	h.Store.MarkStepCompleted(req.StepID)
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
	stepNum, _ := strconv.Atoi(strings.TrimPrefix(stepID, "step"))
	if !exists || stepNum != h.Store.Get() {
		json.NewEncoder(w).Encode(map[string]any{"eligible": false})
		return
	}

	startedAt := h.Store.EnsureStepStarted(stepID)
	elapsed := time.Since(time.Unix(startedAt, 0))

	if elapsed < hintDelay {
		remaining := int((hintDelay - elapsed).Seconds())
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
