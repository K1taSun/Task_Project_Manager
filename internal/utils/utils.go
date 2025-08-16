package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Projekt_go/internal/models"
)

// Join tags into string
func JoinTags(tags []string) string {
	return strings.Join(tags, ",")
}

// Check if string is empty
func IsEmptyString(s string) bool {
	return s == ""
}

// Check if string has appropriate length
func CheckStringLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// Write JSON error
func WriteJSONError(w http.ResponseWriter, status int, msg string) {
	log.Printf("HTTP Error %d: %s", status, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Write JSON message
func WriteJSONMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

// Logging middleware
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// Check if project exists
func ProjectExists(id int) bool {
	return models.ProjectExistsByID(id)
}

// Check if task exists
func TaskExists(id int) bool {
	return models.TaskExistsByID(id)
}

// CORS middleware
func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

// Parse ID from URL path
func ParseIDFromPath(path, prefix string) (int, error) {
	idStr := path[len(prefix):]
	return strconv.Atoi(idStr)
}

// Format time for display
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Parse time from string
func ParseTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}

	layouts := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid time format")
}
