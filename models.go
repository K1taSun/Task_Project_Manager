package main

import (
	"sync"
	"time"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	ID        int        `json:"id"`
	ProjectID int        `json:"project_id"`
	Title     string     `json:"title"`
	Deadline  *time.Time `json:"deadline,omitempty"` // Zmiana na pointer aby obsłużyć null
	Tags      []string   `json:"tags"`
	Priority  int        `json:"priority"`
	Done      bool       `json:"done"`
}

var (
	projects      = make(map[int]Project)
	tasks         = make(map[int]Task)
	mutex         sync.RWMutex // Zmiana na RWMutex dla lepszej wydajności
	nextProjectID = 1
	nextTaskID    = 1
)

// Thread-safe generowanie ID
func generateProjectID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextProjectID
	nextProjectID++
	return id
}

func generateTaskID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextTaskID
	nextTaskID++
	return id
}
