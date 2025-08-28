package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"Projekt_go/internal/models"
	"Projekt_go/internal/storage"
	"Projekt_go/internal/utils"
	"Projekt_go/internal/validation"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Projects handler - GET and POST
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		projects := models.GetAllProjects()

		// Return empty array if no projects
		if projects == nil {
			projects = []models.Project{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error encoding projects")
			return
		}

	case http.MethodPost:
		var p models.Project
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Set default values
		now := time.Now()
		p.CreatedAt = now
		p.UpdatedAt = now
		if p.Status == "" {
			p.Status = "active"
		}
		if p.Badge == nil {
			p.Badge = models.CreateDefaultProjectBadge()
		}

		err = validation.ValidateProject(&p)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		p.ID = models.GenerateProjectID()
		models.SaveProject(p)

		err = storage.SaveProjects()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving projects")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error encoding project")
			return
		}

	default:
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Single project handler - GET, PUT, DELETE
func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDFromPath(r.URL.Path, "/projects/")
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Validate ID
	err = validation.ValidateID(id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	project, exists := models.GetProjectByID(id)
	if !exists {
		utils.WriteJSONError(w, http.StatusNotFound, "Project not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(project)

	case http.MethodPut:
		var updatedProject models.Project
		err := json.NewDecoder(r.Body).Decode(&updatedProject)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Preserve original ID and creation date
		updatedProject.ID = project.ID
		updatedProject.CreatedAt = project.CreatedAt
		models.UpdateProjectTimestamp(&updatedProject)

		err = validation.ValidateProject(&updatedProject)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		models.SaveProject(updatedProject)

		err = storage.SaveProjects()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving projects")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedProject)

	case http.MethodDelete:
		success := models.DeleteProject(id)
		if !success {
			utils.WriteJSONError(w, http.StatusNotFound, "Project not found")
			return
		}

		err = storage.SaveProjects()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving projects")
			return
		}

		broadcastChange()
		utils.WriteJSONMessage(w, http.StatusOK, "Project deleted successfully")

	default:
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Tasks handler - GET and POST
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := models.GetAllTasks()

		// Return empty array if no tasks
		if tasks == nil {
			tasks = []models.Task{}
		}

		// Filter tasks
		tag := r.URL.Query().Get("tag")
		minPriority, _ := strconv.Atoi(r.URL.Query().Get("min_priority"))
		maxPriority, _ := strconv.Atoi(r.URL.Query().Get("max_priority"))
		before := r.URL.Query().Get("before")
		after := r.URL.Query().Get("after")
		projectID, _ := strconv.Atoi(r.URL.Query().Get("project_id"))

		var beforeTime, afterTime time.Time
		if before != "" {
			beforeTime, _ = time.Parse(time.RFC3339, before)
		}
		if after != "" {
			afterTime, _ = time.Parse(time.RFC3339, after)
		}

		var filtered []models.Task
		for _, t := range tasks {
			// Check tag
			if tag != "" {
				found := false
				for _, taskTag := range t.Tags {
					if taskTag == tag {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Check priority range
			if minPriority > 0 && t.Priority < minPriority {
				continue
			}
			if maxPriority > 0 && t.Priority > maxPriority {
				continue
			}

			// Check date range
			if !beforeTime.IsZero() && t.CreatedAt.After(beforeTime) {
				continue
			}
			if !afterTime.IsZero() && t.CreatedAt.Before(afterTime) {
				continue
			}

			// Check project ID
			if projectID > 0 && t.ProjectID != projectID {
				continue
			}

			filtered = append(filtered, t)
		}

		// Ensure we always return an array, even if empty
		if filtered == nil {
			filtered = []models.Task{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(filtered); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error encoding tasks")
			return
		}

	case http.MethodPost:
		var t models.Task
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Set default values
		now := time.Now()
		t.CreatedAt = now
		t.UpdatedAt = now
		if t.Priority == 0 {
			t.Priority = 1
		}
		if t.Badge == nil {
			t.Badge = models.CreatePriorityBadge(t.Priority)
		}

		// Validate project exists
		if !models.ProjectExistsByID(t.ProjectID) {
			utils.WriteJSONError(w, http.StatusBadRequest, "Project not found")
			return
		}

		err = validation.ValidateTask(&t)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		t.ID = models.GenerateTaskID()
		models.SaveTask(t)

		err = storage.SaveTasks()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving tasks")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	default:
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Single task handler - GET, PUT, DELETE
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Validate ID
	err = validation.ValidateID(id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, exists := models.GetTaskByID(id)
	if !exists {
		utils.WriteJSONError(w, http.StatusNotFound, "Task not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodPut:
		var updatedTask models.Task
		err := json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Preserve original ID and creation date
		updatedTask.ID = task.ID
		updatedTask.CreatedAt = task.CreatedAt
		models.UpdateTaskTimestamp(&updatedTask)

		// Validate project exists
		if !models.ProjectExistsByID(updatedTask.ProjectID) {
			utils.WriteJSONError(w, http.StatusBadRequest, "Project not found")
			return
		}

		err = validation.ValidateTask(&updatedTask)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		models.SaveTask(updatedTask)

		err = storage.SaveTasks()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving tasks")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTask)

	case http.MethodDelete:
		success := models.DeleteTask(id)
		if !success {
			utils.WriteJSONError(w, http.StatusNotFound, "Task not found")
			return
		}

		err = storage.SaveTasks()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving tasks")
			return
		}

		broadcastChange()
		utils.WriteJSONMessage(w, http.StatusOK, "Task deleted successfully")

	default:
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Export handler
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		data, err := storage.ExportToJSON()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error exporting data")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=export.json")
		if _, err := w.Write(data); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error writing JSON data")
			return
		}

	case "csv":
		csv, err := storage.ExportToCSV()
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error exporting data")
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
		if _, err := w.Write([]byte(csv)); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error writing CSV data")
			return
		}

	default:
		utils.WriteJSONError(w, http.StatusBadRequest, "Unsupported format")
	}
}

// Stats handler
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tasks := models.GetAllTasks()
	completedTasks := 0
	totalEstimatedHours := 0.0
	totalActualHours := 0.0

	for _, task := range tasks {
		if task.Done {
			completedTasks++
		}
		totalEstimatedHours += task.EstimatedHours
		totalActualHours += task.ActualHours
	}

	totalTasks := models.GetTasksCount()
	var completionRate float64
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	var efficiency float64
	if totalActualHours > 0 && totalEstimatedHours > 0 {
		efficiency = (totalEstimatedHours / totalActualHours) * 100
	}

	stats := map[string]interface{}{
		"total_projects":        models.GetProjectsCount(),
		"total_tasks":           totalTasks,
		"completed_tasks":       completedTasks,
		"pending_tasks":         totalTasks - completedTasks,
		"completion_rate":       completionRate,
		"total_estimated_hours": totalEstimatedHours,
		"total_actual_hours":    totalActualHours,
		"efficiency":            efficiency,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error encoding stats")
		return
	}
}

// Badges handler
func BadgesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	badges := []models.Badge{
		{Text: "New", Color: "#ffffff", Background: "#3b82f6", Icon: "star", Type: "status"},
		{Text: "In Progress", Color: "#ffffff", Background: "#f59e0b", Icon: "hourglass_empty", Type: "status"},
		{Text: "Completed", Color: "#ffffff", Background: "#10b981", Icon: "check_circle", Type: "status"},
		{Text: "Archived", Color: "#ffffff", Background: "#6b7280", Icon: "archive", Type: "status"},
		{Text: "Low", Color: "#ffffff", Background: "#10b981", Icon: "arrow_downward", Type: "priority"},
		{Text: "Medium", Color: "#ffffff", Background: "#f59e0b", Icon: "remove", Type: "priority"},
		{Text: "High", Color: "#ffffff", Background: "#ef4444", Icon: "arrow_upward", Type: "priority"},
		{Text: "Critical", Color: "#ffffff", Background: "#dc2626", Icon: "priority_high", Type: "priority"},
		{Text: "Urgent", Color: "#ffffff", Background: "#7c2d12", Icon: "emergency", Type: "priority"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(badges)
}

// WebSocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Could not upgrade connection")
		return
	}
	defer conn.Close()

	// Send initial data
	data := map[string]interface{}{
		"type":     "initial",
		"projects": models.GetAllProjects(),
		"tasks":    models.GetAllTasks(),
	}

	err = conn.WriteJSON(data)
	if err != nil {
		return
	}

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Broadcast changes to WebSocket clients
func broadcastChange() {
	// This would be implemented with a proper WebSocket manager
	// For now, it's a placeholder
}

// Setup routes
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/projects", utils.WithCORS(utils.LogMiddleware(ProjectsHandler)))
	mux.HandleFunc("/projects/", utils.WithCORS(utils.LogMiddleware(ProjectHandler)))
	mux.HandleFunc("/tasks", utils.WithCORS(utils.LogMiddleware(TasksHandler)))
	mux.HandleFunc("/tasks/", utils.WithCORS(utils.LogMiddleware(TaskHandler)))
	mux.HandleFunc("/export", utils.WithCORS(utils.LogMiddleware(ExportHandler)))
	mux.HandleFunc("/stats", utils.WithCORS(utils.LogMiddleware(StatsHandler)))
	mux.HandleFunc("/badges", utils.WithCORS(utils.LogMiddleware(BadgesHandler)))
	mux.HandleFunc("/ws", utils.LogMiddleware(WSHandler))

	// Static files
	mux.HandleFunc("/", utils.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/static/index.html")
		} else {
			http.NotFound(w, r)
		}
	}))

	return mux
}
