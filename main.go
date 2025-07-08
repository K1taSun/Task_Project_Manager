package main

import (
	"log"
	"net/http"
)

func main() {
	loadData()
	// Rejestracja endpointów
	http.HandleFunc("/projects", withCORS(projectsHandler))
	http.HandleFunc("/projects/", withCORS(projectHandler))
	// http.HandleFunc("/projects/", withCORS(projectTasksHandler)) // Usuwanie powielonej rejestracji
	http.HandleFunc("/tasks", withCORS(tasksHandler))
	http.HandleFunc("/tasks/", withCORS(taskHandler))
	http.HandleFunc("/export", withCORS(exportHandler))
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
