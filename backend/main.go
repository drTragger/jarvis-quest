package main

import (
	"encoding/json"
	"jarvis-quest/backend/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/answer", handlers.ValidateAnswer)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	fs := http.FileServer(http.Dir("./dist"))
	mux.Handle("/", fs)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
