package main

import (
	"encoding/json"
	"log"
	"os"
)

// wczytuje projekty z pliku
func LoadProjects() error {
	file, err := os.Open(ProjectsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik projektow nie istnieje, tworze nowy")
			return SaveProjects()
		}
		return err
	}
	defer file.Close()

	var list []Project
	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()
	projects = make(map[int]Project)
	maxID := 0
	if list != nil {
		for _, p := range list {
			projects[p.ID] = p
			if p.ID > maxID {
				maxID = p.ID
			}
		}
	}
	nextProjectID = maxID + 1
	return nil
}

// zapisuje projekty do pliku
func SaveProjects() error {
	mutex.RLock()
	var list []Project
	for _, p := range projects {
		list = append(list, p)
	}
	mutex.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ProjectsFile, data, 0644)
}

// wczytuje zadania z pliku
func LoadTasks() error {
	file, err := os.Open(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik zadan nie istnieje, tworze nowy")
			return SaveTasks()
		}
		return err
	}
	defer file.Close()

	var list []Task
	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()
	tasks = make(map[int]Task)
	maxID := 0
	if list != nil {
		for _, t := range list {
			tasks[t.ID] = t
			if t.ID > maxID {
				maxID = t.ID
			}
		}
	}
	nextTaskID = maxID + 1
	return nil
}

// zapisuje zadania do pliku
func SaveTasks() error {
	mutex.RLock()
	var list []Task
	for _, p := range tasks {
		list = append(list, p)
	}
	mutex.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(TasksFile, data, 0644)
}
