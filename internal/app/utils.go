package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func joinTags(tags []string) string {
	return strings.Join(tags, ",")
}

func isEmptyString(s string) bool {
	return s == ""
}

func checkStringLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	log.Printf("HTTP Error %d: %s", status, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSONMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func projectExists(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}

// funkcja do powiadamiania o zmianach - deklaracja w handlers.go
