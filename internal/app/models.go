package app

import (
	"sync"
	"time"
)

// Project reprezentuje projekt biznesowy.
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Badge       *Badge    `json:"badge,omitempty"`
	Status      string    `json:"status"`
	Color       string    `json:"color,omitempty"`
}

// Task reprezentuje zadanie z przypisaniem do projektu.
type Task struct {
	ID             int        `json:"id"`
	ProjectID      int        `json:"project_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	Deadline       *time.Time `json:"deadline,omitempty"`
	Tags           []string   `json:"tags"`
	Priority       int        `json:"priority"`
	Done           bool       `json:"done"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Badge          *Badge     `json:"badge,omitempty"`
	Assignee       string     `json:"assignee,omitempty"`
	EstimatedHours float64    `json:"estimated_hours,omitempty"`
	ActualHours    float64    `json:"actual_hours,omitempty"`
}

// Badge opisuje odznakę przypisaną do projektu lub zadania.
type Badge struct {
	Text       string `json:"text"`
	Color      string `json:"color"`
	Background string `json:"background"`
	Icon       string `json:"icon,omitempty"`
	Type       string `json:"type"`
}

var (
	projects      = make(map[int]Project)
	tasks         = make(map[int]Task)
	mutex         sync.RWMutex
	nextProjectID = 1
	nextTaskID    = 1
)

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

func projectExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}

func taskExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := tasks[id]
	return exists
}

func getProjectsCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(projects)
}

func getTasksCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(tasks)
}

func resetIDCounters() {
	mutex.Lock()
	defer mutex.Unlock()
	nextProjectID = 1
	nextTaskID = 1
}

func hasProjects() bool {
	return getProjectsCount() > 0
}

func hasTasks() bool {
	return getTasksCount() > 0
}

func createDefaultProjectBadge() *Badge {
	return &Badge{
		Text:       "Nowy",
		Color:      "#ffffff",
		Background: "#3b82f6",
		Icon:       "star",
		Type:       "status",
	}
}

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

func updateProjectTimestamp(p *Project) {
	p.UpdatedAt = time.Now()
}

func updateTaskTimestamp(t *Task) {
	t.UpdatedAt = time.Now()
}
