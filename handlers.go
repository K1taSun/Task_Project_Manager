package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// funkcja do CORS
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

// obsługa projektów - GET i POST
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Project
		for _, p := range projects {
			list = append(list, p)
		}
		mutex.RUnlock()

		// jak pusta lista to zwracamy pustą tablicę
		if list == nil {
			list = []Project{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var p Project
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		
		// ustawiamy domyślne wartości
		now := time.Now()
		p.CreatedAt = now
		p.UpdatedAt = now
		if p.Status == "" {
			p.Status = "active"
		}
		if p.Badge == nil {
			p.Badge = createDefaultProjectBadge()
		}
		
		err = validateProject(&p)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		p.ID = generateProjectID()
		mutex.Lock()
		projects[p.ID] = p
		mutex.Unlock()
		err = SaveProjects()
		if err != nil {
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

// obsługa pojedynczego projektu - GET, PUT, DELETE
func projectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/projects/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// sprawdzamy czy ID jest poprawne
	err = validateID(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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
		err := json.NewDecoder(r.Body).Decode(&updatedProject)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		
		// zachowujemy oryginalne ID i datę utworzenia
		updatedProject.ID = p.ID
		updatedProject.CreatedAt = p.CreatedAt
		updatedProject.UpdatedAt = time.Now()
		
		err = validateProject(&updatedProject)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		
		mutex.Lock()
		projects[id] = updatedProject
		mutex.Unlock()
		
		err = SaveProjects()
		if err != nil {
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
		
		err = SaveProjects()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Blad zapisu projektow")
			return
		}
		
		broadcastChange()
		writeJSONMessage(w, http.StatusOK, "Project deleted successfully")
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// obsługa zadań - GET i POST
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Task
		for _, t := range tasks {
			list = append(list, t)
		}
		mutex.RUnlock()

		// jak pusta lista to zwracamy pustą tablicę
		if list == nil {
			list = []Task{}
		}

		// filtrowanie zadań
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
			// sprawdzamy tag
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
			// sprawdzamy priorytet
			if minPriority != 0 && t.Priority < minPriority {
				continue
			}
			if maxPriority != 0 && t.Priority > maxPriority {
				continue
			}
			// sprawdzamy deadline
			if !beforeTime.IsZero() && t.Deadline != nil && t.Deadline.After(beforeTime) {
				continue
			}
			if !afterTime.IsZero() && t.Deadline != nil && t.Deadline.Before(afterTime) {
				continue
			}
			filtered = append(filtered, t)
		}

		// sortowanie zadań
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
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		
		// ustawiamy domyślne wartości
		now := time.Now()
		t.CreatedAt = now
		t.UpdatedAt = now
		if t.Badge == nil {
			t.Badge = createPriorityBadge(t.Priority)
		}
		
		if err := validateTask(&t); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		t.ID = generateTaskID()
		mutex.Lock()
		tasks[t.ID] = t
		mutex.Unlock()
		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}
		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /projects/{id}/tasks (GET, POST)
func projectTasksHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) < len("/projects/") {
		writeJSONError(w, http.StatusNotFound, "Not found")
		return
	}
	trimmed := path[len("/projects/"):] // np. 123/tasks lub 123
	parts := []rune(trimmed)
	idStr := ""
	for i, c := range parts {
		if c < '0' || c > '9' {
			idStr = string(parts[:i])
			trimmed = string(parts[i:])
			break
		}
	}
	if idStr == "" {
		idStr = trimmed
		trimmed = ""
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}
	if trimmed != "/tasks" {
		writeJSONError(w, http.StatusNotFound, "Not found")
		return
	}

	// Sprawdź czy projekt istnieje
	mutex.RLock()
	_, projectExists := projects[id]
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
			if t.ProjectID == id {
				list = append(list, t)
			}
		}
		mutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if err := validateTask(&t); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		t.ID = generateTaskID()
		t.ProjectID = id
		mutex.Lock()
		tasks[t.ID] = t
		mutex.Unlock()
		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}
		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func taskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
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
		var updated Task
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if err := validateTask(&updated); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		updated.ID = id
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

// /export?format=csv|json
func exportHandler(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "csv" {
		exportCSV(w, r)
	} else {
		exportJSON(w, r)
	}
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
	// Projekty
	csvWriter.Write([]string{"ProjectID", "ProjectName"})
	for _, p := range projects {
		csvWriter.Write([]string{strconv.Itoa(p.ID), p.Name})
	}
	csvWriter.Write([]string{})
	// Zadania
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

// /badges - zarządzanie odznakami
func badgesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// zwracamy dostępne typy odznak
		badgeTypes := map[string]interface{}{
			"priority": map[string]interface{}{
				"1": createPriorityBadge(1),
				"2": createPriorityBadge(2),
				"3": createPriorityBadge(3),
				"4": createPriorityBadge(4),
				"5": createPriorityBadge(5),
			},
			"status": map[string]interface{}{
				"active": &Badge{Text: "Aktywny", Color: "#ffffff", Background: "#10b981", Icon: "play_circle", Type: "status"},
				"completed": &Badge{Text: "Ukończony", Color: "#ffffff", Background: "#3b82f6", Icon: "check_circle", Type: "status"},
				"archived": &Badge{Text: "Zarchiwizowany", Color: "#ffffff", Background: "#6b7280", Icon: "archive", Type: "status"},
			},
			"category": map[string]interface{}{
				"bug": &Badge{Text: "Błąd", Color: "#ffffff", Background: "#ef4444", Icon: "bug_report", Type: "category"},
				"feature": &Badge{Text: "Funkcja", Color: "#ffffff", Background: "#8b5cf6", Icon: "new_releases", Type: "category"},
				"improvement": &Badge{Text: "Ulepszenie", Color: "#ffffff", Background: "#f59e0b", Icon: "trending_up", Type: "category"},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(badgeTypes)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /stats - statystyki aplikacji
func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	mutex.RLock()
	projectCount := len(projects)
	taskCount := len(tasks)
	
	// liczymy zadania według priorytetu
	priorityStats := make(map[int]int)
	for _, task := range tasks {
		priorityStats[task.Priority]++
	}
	
	// liczymy ukończone zadania
	completedTasks := 0
	for _, task := range tasks {
		if task.Done {
			completedTasks++
		}
	}
	
	// liczymy zadania według projektu
	projectTaskStats := make(map[int]int)
	for _, task := range tasks {
		projectTaskStats[task.ProjectID]++
	}
	
	stats := map[string]interface{}{
		"projects": map[string]interface{}{
			"total": projectCount,
			"active": func() int {
				count := 0
				for _, p := range projects {
					if p.Status == "active" {
						count++
					}
				}
				return count
			}(),
		},
		"tasks": map[string]interface{}{
			"total": taskCount,
			"completed": completedTasks,
			"pending": taskCount - completedTasks,
			"by_priority": priorityStats,
			"by_project": projectTaskStats,
		},
		"system": map[string]interface{}{
			"uptime": time.Since(startTime).String(),
			"version": "2.0.0",
		},
	}
	mutex.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// WebSocket upgrade handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // pozwalamy na wszystkie originy w trybie deweloperskim
		},
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	
	// dodajemy połączenie do listy aktywnych
	wsConnections = append(wsConnections, conn)
	defer func() {
		// usuwamy połączenie z listy
		for i, c := range wsConnections {
			if c == conn {
				wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
				break
			}
		}
	}()
	
	// nasłuchujemy na wiadomości
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		
		// echo message back (można rozszerzyć o obsługę komend)
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// globalne zmienne dla WebSocket
var (
	wsConnections []*websocket.Conn
	startTime     = time.Now()
)

// funkcja do powiadamiania o zmianach przez WebSocket
func broadcastChange() {
	if len(wsConnections) == 0 {
		return
	}
	
	message := map[string]string{
		"type": "change",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}
	
	// wysyłamy wiadomość do wszystkich aktywnych połączeń
	for _, conn := range wsConnections {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
		}
	}
}
