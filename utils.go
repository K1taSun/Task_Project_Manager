package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// łączy tagi w string
// TODO: zrobić to lepiej z strings.Join
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

// sprawdza czy string jest pusty (prosta funkcja)
func isEmptyString(s string) bool {
	if s == "" {
		return true
	}
	return false
}

// sprawdza czy string ma odpowiednią długość
func checkStringLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// pisze błąd JSON
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	log.Printf("HTTP Error %d: %s", status, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// pisze wiadomość JSON
func writeJSONMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

// middleware do logowania
// TODO: dodać więcej informacji w logach
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// sprawdza czy projekt istnieje
func projectExists(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}

// funkcja do powiadamiania o zmianach (pusta na razie)
// TODO: dodać powiadomienia o zmianach
// na razie nic nie robi

// prosta funkcja do sprawdzania czy request jest GET
func isGetRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// prosta funkcja do sprawdzania czy request jest POST
func isPostRequest(r *http.Request) bool {
	return r.Method == http.MethodPost
}
