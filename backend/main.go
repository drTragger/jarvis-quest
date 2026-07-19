package main

import (
	"encoding/json"
	"jarvis-quest/backend/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	mux.HandleFunc("/api/onboarding/seen", quest.MarkOnboardingSeen)
	mux.HandleFunc("/api/rules/seen", quest.MarkRulesSeen)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.Handle("/", spaHandler("./dist"))

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// spaHandler serves static files from root, falling back to index.html for
// any path that doesn't match a real file so client-side routing (history
// mode, no #) works on direct loads and refreshes.
func spaHandler(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	indexPath := filepath.Join(root, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean(r.URL.Path)
		fullPath := filepath.Join(root, cleanPath)

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, indexPath)
			return
		}
		fs.ServeHTTP(w, r)
	})
}
