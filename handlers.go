package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

// /projects (GET, POST)
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Project
		for _, p := range projects {
			list = append(list, p)
		}
		mutex.RUnlock()

		// Jeśli lista jest pusta, zwróć pustą tablicę zamiast null
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
		if err := validateProject(&p); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		p.ID = generateProjectID()
		mutex.Lock()
		projects[p.ID] = p
		mutex.Unlock()
		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu projektów")
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

// /projects/{id} (GET, PUT, DELETE)
func projectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/projects/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid project ID")
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
		var updated Project
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if err := validateProject(&updated); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		updated.ID = id
		mutex.Lock()
		projects[id] = updated
		mutex.Unlock()
		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu projektów")
			return
		}
		broadcastChange()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		// Usuń wszystkie zadania powiązane z projektem
		mutex.Lock()
		delete(projects, id)
		// Usuń zadania powiązane z projektem
		for taskID, task := range tasks {
			if task.ProjectID == id {
				delete(tasks, taskID)
			}
		}
		mutex.Unlock()
		if err := SaveProjects(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu projektów")
			return
		}
		if err := SaveTasks(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Błąd zapisu zadań")
			return
		}
		broadcastChange()
		writeJSONMessage(w, http.StatusNoContent, "Project deleted")
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /tasks (GET, POST)
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		var list []Task
		for _, t := range tasks {
			list = append(list, t)
		}
		mutex.RUnlock()

		// Jeśli lista jest pusta, zwróć pustą tablicę zamiast null
		if list == nil {
			list = []Task{}
		}

		// Filtrowanie
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

		// Sortowanie
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
