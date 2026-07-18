package handlers

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type ProgressStore struct {
	mu             sync.Mutex
	path           string
	step           int
	attempts       map[string]int
	startedAt      map[string]int64
	completedAt    map[string]int64
	onboardingSeen bool
	rulesSeen      bool
}

type persistedState struct {
	Step           int              `json:"step"`
	Attempts       map[string]int   `json:"attempts"`
	StartedAt      map[string]int64 `json:"started_at"`
	CompletedAt    map[string]int64 `json:"completed_at"`
	OnboardingSeen bool             `json:"onboarding_seen"`
	RulesSeen      bool             `json:"rules_seen"`
}

func NewProgressStore(path string) *ProgressStore {
	s := &ProgressStore{
		path:        path,
		step:        1,
		attempts:    make(map[string]int),
		startedAt:   make(map[string]int64),
		completedAt: make(map[string]int64),
	}
	s.load()
	return s
}

func (s *ProgressStore) load() {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var data persistedState
	if json.Unmarshal(b, &data) != nil {
		return
	}
	if data.Step > 0 {
		s.step = data.Step
	}
	if data.Attempts != nil {
		s.attempts = data.Attempts
	}
	if data.StartedAt != nil {
		s.startedAt = data.StartedAt
	}
	if data.CompletedAt != nil {
		s.completedAt = data.CompletedAt
	}
	s.onboardingSeen = data.OnboardingSeen
	s.rulesSeen = data.RulesSeen
}

func (s *ProgressStore) save() {
	data := persistedState{
		Step:           s.step,
		Attempts:       s.attempts,
		StartedAt:      s.startedAt,
		CompletedAt:    s.completedAt,
		OnboardingSeen: s.onboardingSeen,
		RulesSeen:      s.rulesSeen,
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(s.path, b, 0644)
}

func (s *ProgressStore) Get() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.step
}

func (s *ProgressStore) Advance(step int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if step > s.step {
		s.step = step
	}
	val := s.step
	s.save()
	return val
}

func (s *ProgressStore) EnsureStepStarted(stepID string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ts, ok := s.startedAt[stepID]; ok {
		return ts
	}
	now := time.Now().Unix()
	s.startedAt[stepID] = now
	s.save()
	return now
}

func (s *ProgressStore) IncrementAttempts(stepID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempts[stepID]++
	val := s.attempts[stepID]
	s.save()
	return val
}

func (s *ProgressStore) GetAttempts(stepID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.attempts[stepID]
}

func (s *ProgressStore) MarkStepCompleted(stepID string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ts, ok := s.completedAt[stepID]; ok {
		return ts
	}
	now := time.Now().Unix()
	s.completedAt[stepID] = now
	s.save()
	return now
}

func (s *ProgressStore) GetOnboardingSeen() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.onboardingSeen
}

func (s *ProgressStore) MarkOnboardingSeen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onboardingSeen = true
	s.save()
}

func (s *ProgressStore) GetRulesSeen() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rulesSeen
}

func (s *ProgressStore) MarkRulesSeen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rulesSeen = true
	s.save()
}
