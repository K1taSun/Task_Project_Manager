package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	projectsFile = "data_projects.json"
	tasksFile    = "data_tasks.json"
)

func LoadProjects() error {
	file, err := os.Open(projectsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik projektów nie istnieje, tworzenie nowego")
			return SaveProjects() // Utwórz pusty plik
		}
		return err
	}
	defer file.Close()

	var list []Project
	if err := json.NewDecoder(file).Decode(&list); err != nil {
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
	return ioutil.WriteFile(projectsFile, data, 0644)
}

func LoadTasks() error {
	file, err := os.Open(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik zadań nie istnieje, tworzenie nowego")
			return SaveTasks() // Utwórz pusty plik
		}
		return err
	}
	defer file.Close()

	var list []Task
	if err := json.NewDecoder(file).Decode(&list); err != nil {
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

func SaveTasks() error {
	mutex.RLock()
	var list []Task
	for _, t := range tasks {
		list = append(list, t)
	}
	mutex.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tasksFile, data, 0644)
}
