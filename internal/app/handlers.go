package app

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RegisterRoutes dodaje wszystkie endpointy HTTP aplikacji do przekazanego muxa.
func RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		mux = http.DefaultServeMux
	}

	mux.HandleFunc("/projects", withCORS(logMiddleware(projectsHandler)))
	mux.HandleFunc("/projects/", withCORS(logMiddleware(projectHandler)))
	mux.HandleFunc("/tasks", withCORS(logMiddleware(tasksHandler)))
	mux.HandleFunc("/tasks/", withCORS(logMiddleware(taskHandler)))
	mux.HandleFunc("/export", withCORS(logMiddleware(exportHandler)))
	mux.HandleFunc("/badges", withCORS(logMiddleware(badgesHandler)))
	mux.HandleFunc("/ws", logMiddleware(wsHandler))
	mux.HandleFunc("/", logMiddleware(serveIndex))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
		return
	}
	http.NotFound(w, r)
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
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

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Project
		for _, p := range projects {
			list = append(list, p)
		}
		mutex.RUnlock()

		if list == nil {
			list = []Project{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var p Project
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		now := time.Now()
		p.CreatedAt = now
		p.UpdatedAt = now
		if p.Status == "" {
			p.Status = "active"
		}
		if p.Badge == nil {
			p.Badge = createDefaultProjectBadge()
		}

		if err := validateProject(&p); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		p.ID = generateProjectID()
		mutex.Lock()
		projects[p.ID] = p
		mutex.Unlock()

		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Blad zapisu projektow")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, "/projects/")
	parts := strings.FieldsFunc(trimmedPath, func(r rune) bool { return r == '/' })
	if len(parts) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := validateID(id); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(parts) > 1 && parts[1] == "tasks" {
		projectTasksHandler(w, r, id)
		return
	}

	mutex.RLock()
	p, ok := projects[id]
	mutex.RUnlock()
	if !ok {
		writeJSONError(w, http.StatusNotFound, "Project not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	case http.MethodPut:
		var updatedProject Project
		if err := json.NewDecoder(r.Body).Decode(&updatedProject); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		updatedProject.ID = p.ID
		updatedProject.CreatedAt = p.CreatedAt
		updatedProject.UpdatedAt = time.Now()

		if err := validateProject(&updatedProject); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		mutex.Lock()
		projects[id] = updatedProject
		mutex.Unlock()

		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Blad zapisu projektow")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedProject)
	case http.MethodDelete:
		mutex.Lock()
		delete(projects, id)
		mutex.Unlock()

		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Blad zapisu projektow")
			return
		}

		broadcastChange()
		writeJSONMessage(w, http.StatusOK, "Project deleted successfully")
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Task
		for _, t := range tasks {
			list = append(list, t)
		}
		mutex.RUnlock()

		if list == nil {
			list = []Task{}
		}

		tag := r.URL.Query().Get("tag")
		minPriority, _ := strconv.Atoi(r.URL.Query().Get("min_priority"))
		maxPriority, _ := strconv.Atoi(r.URL.Query().Get("max_priority"))
		before := r.URL.Query().Get("before")
		after := r.URL.Query().Get("after")

		var beforeTime, afterTime time.Time
		if before != "" {
			beforeTime, _ = time.Parse(time.RFC3339, before)
		}
		if after != "" {
			afterTime, _ = time.Parse(time.RFC3339, after)
		}

		var filtered []Task
		for _, t := range list {
			if tag != "" {
				found := false
				for _, tg := range t.Tags {
					if strings.EqualFold(tg, tag) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if minPriority != 0 && t.Priority < minPriority {
				continue
			}
			if maxPriority != 0 && t.Priority > maxPriority {
				continue
			}
			if !beforeTime.IsZero() && t.Deadline != nil && t.Deadline.After(beforeTime) {
				continue
			}
			if !afterTime.IsZero() && t.Deadline != nil && t.Deadline.Before(afterTime) {
				continue
			}
			filtered = append(filtered, t)
		}

		sortBy := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")
		if sortBy != "" {
			sort.Slice(filtered, func(i, j int) bool {
				switch sortBy {
				case "priority":
					if order == "desc" {
						return filtered[i].Priority > filtered[j].Priority
					}
					return filtered[i].Priority < filtered[j].Priority
				case "deadline":
					if order == "desc" {
						if filtered[i].Deadline == nil {
							return false
						}
						if filtered[j].Deadline == nil {
							return true
						}
						return filtered[i].Deadline.After(*filtered[j].Deadline)
					}
					if filtered[i].Deadline == nil {
						return false
					}
					if filtered[j].Deadline == nil {
						return true
					}
					return filtered[i].Deadline.Before(*filtered[j].Deadline)
				case "title":
					if order == "desc" {
						return filtered[i].Title > filtered[j].Title
					}
					return filtered[i].Title < filtered[j].Title
				}
				return false
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if filtered == nil {
			filtered = []Task{}
		}
		json.NewEncoder(w).Encode(filtered)
	case http.MethodPost:
		var payload taskPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		now := time.Now()
		base := applyTaskPayload(Task{}, payload)
		base.CreatedAt = now
		base.UpdatedAt = now

		if base.Badge == nil {
			base.Badge = createPriorityBadge(base.Priority)
		}

		if err := validateTask(&base); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if base.ProjectID != 0 {
			if err := validateID(base.ProjectID); err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
			if !projectExistsByID(base.ProjectID) {
				writeJSONError(w, http.StatusBadRequest, "Project does not exist")
				return
			}
		}

		base.ID = generateTaskID()
		mutex.Lock()
		tasks[base.ID] = base
		mutex.Unlock()

		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(base)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func projectTasksHandler(w http.ResponseWriter, r *http.Request, projectID int) {
	mutex.RLock()
	_, projectExists := projects[projectID]
	mutex.RUnlock()
	if !projectExists {
		writeJSONError(w, http.StatusNotFound, "Project not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Task
		for _, t := range tasks {
			if t.ProjectID == projectID {
				list = append(list, t)
			}
		}
		mutex.RUnlock()
		if list == nil {
			list = []Task{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var payload taskPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		now := time.Now()
		base := applyTaskPayload(Task{}, payload)
		base.ProjectID = projectID
		base.CreatedAt = now
		base.UpdatedAt = now

		if base.Badge == nil {
			base.Badge = createPriorityBadge(base.Priority)
		}

		if err := validateTask(&base); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		base.ID = generateTaskID()
		mutex.Lock()
		tasks[base.ID] = base
		mutex.Unlock()

		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(base)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	mutex.RLock()
	t, ok := tasks[id]
	mutex.RUnlock()
	if !ok {
		writeJSONError(w, http.StatusNotFound, "Task not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	case http.MethodPut:
		var payload taskPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		updated := applyTaskPayload(t, payload)
		if payload.Priority == nil {
			updated.Priority = t.Priority
		}
		if payload.Done == nil {
			updated.Done = t.Done
		}
		if payload.Deadline == nil {
			updated.Deadline = t.Deadline
		}
		if payload.ProjectID == nil {
			updated.ProjectID = t.ProjectID
		}
		if payload.EstimatedHours == nil {
			updated.EstimatedHours = t.EstimatedHours
		}
		if payload.ActualHours == nil {
			updated.ActualHours = t.ActualHours
		}
		if payload.Badge == nil {
			updated.Badge = t.Badge
		}

		if updated.ProjectID != 0 {
			if err := validateID(updated.ProjectID); err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
			if !projectExistsByID(updated.ProjectID) {
				writeJSONError(w, http.StatusBadRequest, "Project does not exist")
				return
			}
		}

		updated.ID = id
		updated.CreatedAt = t.CreatedAt
		updated.UpdatedAt = time.Now()
		if updated.Badge == nil {
			updated.Badge = createPriorityBadge(updated.Priority)
		}

		if err := validateTask(&updated); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		mutex.Lock()
		tasks[id] = updated
		mutex.Unlock()

		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}

		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		mutex.Lock()
		delete(tasks, id)
		mutex.Unlock()

		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}

		broadcastChange()
		writeJSONMessage(w, http.StatusNoContent, "Task deleted")
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "csv" {
		exportCSV(w, r)
		return
	}
	exportJSON(w, r)
}

func exportJSON(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	data := struct {
		Projects []Project `json:"projects"`
		Tasks    []Task    `json:"tasks"`
	}{}
	for _, p := range projects {
		data.Projects = append(data.Projects, p)
	}
	for _, t := range tasks {
		data.Tasks = append(data.Tasks, t)
	}
	mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func exportCSV(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=export.csv")
	csvWriter := csv.NewWriter(w)

	csvWriter.Write([]string{"ProjectID", "ProjectName"})
	for _, p := range projects {
		csvWriter.Write([]string{strconv.Itoa(p.ID), p.Name})
	}
	csvWriter.Write([]string{})

	csvWriter.Write([]string{"TaskID", "ProjectID", "Title", "Description", "Deadline", "Tags", "Priority", "Done"})
	for _, t := range tasks {
		deadlineStr := ""
		if t.Deadline != nil {
			deadlineStr = t.Deadline.Format(time.RFC3339)
		}
		csvWriter.Write([]string{
			strconv.Itoa(t.ID),
			strconv.Itoa(t.ProjectID),
			t.Title,
			t.Description,
			deadlineStr,
			"[" + joinTags(t.Tags) + "]",
			strconv.Itoa(t.Priority),
			strconv.FormatBool(t.Done),
		})
	}
	mutex.RUnlock()
	csvWriter.Flush()
}

func badgesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		badgeTypes := map[string]interface{}{
			"priority": map[string]interface{}{
				"1": createPriorityBadge(1),
				"2": createPriorityBadge(2),
				"3": createPriorityBadge(3),
				"4": createPriorityBadge(4),
				"5": createPriorityBadge(5),
			},
			"status": map[string]interface{}{
				"active":    &Badge{Text: "Aktywny", Color: "#ffffff", Background: "#10b981", Icon: "play_circle", Type: "status"},
				"completed": &Badge{Text: "Ukończony", Color: "#ffffff", Background: "#3b82f6", Icon: "check_circle", Type: "status"},
				"archived":  &Badge{Text: "Zarchiwizowany", Color: "#ffffff", Background: "#6b7280", Icon: "archive", Type: "status"},
			},
			"category": map[string]interface{}{
				"bug":         &Badge{Text: "Błąd", Color: "#ffffff", Background: "#ef4444", Icon: "bug_report", Type: "category"},
				"feature":     &Badge{Text: "Funkcja", Color: "#ffffff", Background: "#8b5cf6", Icon: "new_releases", Type: "category"},
				"improvement": &Badge{Text: "Ulepszenie", Color: "#ffffff", Background: "#f59e0b", Icon: "trending_up", Type: "category"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(badgeTypes)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	wsMutex.Lock()
	wsConnections = append(wsConnections, conn)
	wsMutex.Unlock()

	defer func() {
		wsMutex.Lock()
		for i, c := range wsConnections {
			if c == conn {
				wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
				break
			}
		}
		wsMutex.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

var (
	wsConnections []*websocket.Conn
	wsMutex       sync.Mutex
)

type taskPayload struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Tags           []string   `json:"tags"`
	Priority       *int       `json:"priority"`
	Done           *bool      `json:"done"`
	Deadline       *time.Time `json:"deadline"`
	ProjectID      *int       `json:"project_id"`
	Assignee       string     `json:"assignee"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ActualHours    *float64   `json:"actual_hours"`
	Badge          *Badge     `json:"badge"`
}

func applyTaskPayload(base Task, payload taskPayload) Task {
	base.Title = payload.Title
	base.Description = payload.Description
	if payload.Tags != nil {
		base.Tags = payload.Tags
	}
	if payload.Priority != nil {
		base.Priority = *payload.Priority
	}
	if payload.Done != nil {
		base.Done = *payload.Done
	}
	if payload.Deadline != nil {
		base.Deadline = payload.Deadline
	}
	if payload.ProjectID != nil {
		base.ProjectID = *payload.ProjectID
	}
	if payload.EstimatedHours != nil {
		base.EstimatedHours = *payload.EstimatedHours
	}
	if payload.ActualHours != nil {
		base.ActualHours = *payload.ActualHours
	}
	if payload.Badge != nil {
		base.Badge = payload.Badge
	}
	base.Assignee = payload.Assignee
	return base
}

func broadcastChange() {
	wsMutex.Lock()
	connections := make([]*websocket.Conn, len(wsConnections))
	copy(connections, wsConnections)
	wsMutex.Unlock()

	if len(connections) == 0 {
		return
	}

	message := map[string]string{
		"type":      "change",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
		}
	}
}
