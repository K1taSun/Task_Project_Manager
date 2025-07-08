package main

import (
	"encoding/json"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSONMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
} 