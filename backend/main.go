package main

import (
	"encoding/json"
	"jarvis-quest/backend/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	if wd, err := os.Getwd(); err == nil {
		log.Printf("Working directory: %s", wd)
	}

	mux := http.NewServeMux()

	store := handlers.NewProgressStore("./progress.json")
	quest := &handlers.QuestHandler{Store: store}

	mux.HandleFunc("/api/answer", quest.ValidateAnswer)
	mux.HandleFunc("/api/progress", quest.GetProgress)
	mux.HandleFunc("/api/hint", quest.GetHint)
	mux.HandleFunc("/api/step/start", quest.StartStep)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	fs := http.FileServer(http.Dir("./dist"))
	mux.Handle("/", fs)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
