package models

import (
	"sync"
	"time"
)

// Project structure
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

// Task structure
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

// Badge structure
type Badge struct {
	Text       string `json:"text"`
	Color      string `json:"color"`
	Background string `json:"background"`
	Icon       string `json:"icon,omitempty"`
	Type       string `json:"type"` // status, priority, category, custom
}

// Global variables
var (
	projects      = make(map[int]Project)
	tasks         = make(map[int]Task)
	mutex         sync.RWMutex
	nextProjectID = 1
	nextTaskID    = 1
)

// Generate ID for project
func GenerateProjectID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextProjectID
	nextProjectID++
	return id
}

// Generate ID for task
func GenerateTaskID() int {
	mutex.Lock()
	defer mutex.Unlock()
	id := nextTaskID
	nextTaskID++
	return id
}

// Check if project exists by ID
func ProjectExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := projects[id]
	return exists
}

// Check if task exists by ID
func TaskExistsByID(id int) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := tasks[id]
	return exists
}

// Get projects count
func GetProjectsCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(projects)
}

// Get tasks count
func GetTasksCount() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(tasks)
}

// Reset ID counters
func ResetIDCounters() {
	mutex.Lock()
	defer mutex.Unlock()
	nextProjectID = 1
	nextTaskID = 1
}

// Reset all data (for testing)
func ResetAllData() {
	mutex.Lock()
	defer mutex.Unlock()
	projects = make(map[int]Project)
	tasks = make(map[int]Task)
	nextProjectID = 1
	nextTaskID = 1
}

// Check if there are any projects
func HasProjects() bool {
	return GetProjectsCount() > 0
}

// Check if there are any tasks
func HasTasks() bool {
	return GetTasksCount() > 0
}

// Create default badge for project
func CreateDefaultProjectBadge() *Badge {
	return &Badge{
		Text:       "New",
		Color:      "#ffffff",
		Background: "#3b82f6",
		Icon:       "star",
		Type:       "status",
	}
}

// Create priority badge for task
func CreatePriorityBadge(priority int) *Badge {
	badges := map[int]*Badge{
		1: {Text: "Low", Color: "#ffffff", Background: "#10b981", Icon: "arrow_downward", Type: "priority"},
		2: {Text: "Medium", Color: "#ffffff", Background: "#f59e0b", Icon: "remove", Type: "priority"},
		3: {Text: "High", Color: "#ffffff", Background: "#ef4444", Icon: "arrow_upward", Type: "priority"},
		4: {Text: "Critical", Color: "#ffffff", Background: "#dc2626", Icon: "priority_high", Type: "priority"},
		5: {Text: "Urgent", Color: "#ffffff", Background: "#7c2d12", Icon: "emergency", Type: "priority"},
	}

	if badge, exists := badges[priority]; exists {
		return badge
	}
	return &Badge{Text: "Unknown", Color: "#ffffff", Background: "#6b7280", Icon: "help", Type: "priority"}
}

// Update project timestamp
func UpdateProjectTimestamp(p *Project) {
	p.UpdatedAt = time.Now()
}

// Update task timestamp
func UpdateTaskTimestamp(t *Task) {
	t.UpdatedAt = time.Now()
}

// Get all projects
func GetAllProjects() []Project {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []Project
	for _, p := range projects {
		list = append(list, p)
	}
	// Always return an empty slice instead of nil
	if list == nil {
		list = []Project{}
	}
	return list
}

// Get all tasks
func GetAllTasks() []Task {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []Task
	for _, t := range tasks {
		list = append(list, t)
	}
	// Always return an empty slice instead of nil
	if list == nil {
		list = []Task{}
	}
	return list
}

// Get project by ID
func GetProjectByID(id int) (Project, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	project, exists := projects[id]
	return project, exists
}

// Get task by ID
func GetTaskByID(id int) (Task, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	task, exists := tasks[id]
	return task, exists
}

// Save project
func SaveProject(project Project) {
	mutex.Lock()
	defer mutex.Unlock()
	projects[project.ID] = project
}

// Save task
func SaveTask(task Task) {
	mutex.Lock()
	defer mutex.Unlock()
	tasks[task.ID] = task
}

// Delete project
func DeleteProject(id int) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := projects[id]; exists {
		delete(projects, id)
		return true
	}
	return false
}

// Delete task
func DeleteTask(id int) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := tasks[id]; exists {
		delete(tasks, id)
		return true
	}
	return false
}

// Get tasks by project ID
func GetTasksByProjectID(projectID int) []Task {
	mutex.RLock()
	defer mutex.RUnlock()
	var projectTasks []Task
	for _, task := range tasks {
		if task.ProjectID == projectID {
			projectTasks = append(projectTasks, task)
		}
	}
	return projectTasks
}
