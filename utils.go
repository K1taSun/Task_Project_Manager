package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func joinTags(tags []string) string {
	returnString := ""
	for i, tag := range tags {
		if i > 0 {
			returnString += ","
		}
		returnString += tag
	}
	return returnString
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

// Middleware do logowania requestów
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// Funkcja do sprawdzania czy projekt istnieje
func projectExists(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}
