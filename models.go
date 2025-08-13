package main

import (
	"sync"
	"time"
)

// struktura projektu
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Badge       *Badge    `json:"badge,omitempty"`
	Status      string    `json:"status"` // active, completed, archived
	Color       string    `json:"color,omitempty"`
}

// struktura zadania
type Task struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Tags        []string   `json:"tags"`
	Priority    int        `json:"priority"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Badge       *Badge     `json:"badge,omitempty"`
	Assignee    string     `json:"assignee,omitempty"`
	EstimatedHours float64 `json:"estimated_hours,omitempty"`
	ActualHours   float64 `json:"actual_hours,omitempty"`
}

// struktura odznaki/badge
type Badge struct {
	Text        string `json:"text"`
	Color       string `json:"color"`
	Background  string `json:"background"`
	Icon        string `json:"icon,omitempty"`
	Type        string `json:"type"` // status, priority, category, custom
}

// globalne zmienne
var (
	projects      = make(map[int]Project)
	tasks         = make(map[int]Task)
	mutex         sync.RWMutex
	nextProjectID = 1
	nextTaskID    = 1
)

// generuje ID dla projektu
func generateProjectID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextProjectID
	nextProjectID++
	return id
}

// generuje ID dla zadania
func generateTaskID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextTaskID
	nextTaskID++
	return id
}

// sprawdza czy projekt istnieje
func projectExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}

// sprawdza czy zadanie istnieje
func taskExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := tasks[id]
	return exists
}

// zwraca liczbę projektów
func getProjectsCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(projects)
}

// zwraca liczbę zadań
func getTasksCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(tasks)
}

// resetuje liczniki ID (prosta funkcja)
func resetIDCounters() {
	mutex.Lock()
	defer mutex.Unlock()
	nextProjectID = 1
	nextTaskID = 1
}

// sprawdza czy są jakieś projekty
func hasProjects() bool {
	return getProjectsCount() > 0
}

// sprawdza czy są jakieś zadania
func hasTasks() bool {
	return getTasksCount() > 0
}

// tworzy domyślną odznakę dla projektu
func createDefaultProjectBadge() *Badge {
	return &Badge{
		Text:       "Nowy",
		Color:      "#ffffff",
		Background: "#3b82f6",
		Icon:       "star",
		Type:       "status",
	}
}

// tworzy odznakę priorytetu dla zadania
func createPriorityBadge(priority int) *Badge {
	badges := map[int]*Badge{
		1: {Text: "Niski", Color: "#ffffff", Background: "#10b981", Icon: "arrow_downward", Type: "priority"},
		2: {Text: "Średni", Color: "#ffffff", Background: "#f59e0b", Icon: "remove", Type: "priority"},
		3: {Text: "Wysoki", Color: "#ffffff", Background: "#ef4444", Icon: "arrow_upward", Type: "priority"},
		4: {Text: "Krytyczny", Color: "#ffffff", Background: "#dc2626", Icon: "priority_high", Type: "priority"},
		5: {Text: "Pilny", Color: "#ffffff", Background: "#7c2d12", Icon: "emergency", Type: "priority"},
	}
	
	if badge, exists := badges[priority]; exists {
		return badge
	}
	return &Badge{Text: "Nieznany", Color: "#ffffff", Background: "#6b7280", Icon: "help", Type: "priority"}
}

// aktualizuje czas modyfikacji projektu
func updateProjectTimestamp(p *Project) {
	p.UpdatedAt = time.Now()
}

// aktualizuje czas modyfikacji zadania
func updateTaskTimestamp(t *Task) {
	t.UpdatedAt = time.Now()
}
