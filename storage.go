package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	projectsFile = "data_projects.json"
	tasksFile    = "data_tasks.json"
)

func loadData() {
	mutex.Lock()
	defer mutex.Unlock()
	// Projects
	if data, err := ioutil.ReadFile(projectsFile); err == nil {
		var list []Project
		if err := json.Unmarshal(data, &list); err == nil {
			for _, p := range list {
				projects[p.ID] = p
			}
		}
	}
	// Tasks
	if data, err := ioutil.ReadFile(tasksFile); err == nil {
		var list []Task
		if err := json.Unmarshal(data, &list); err == nil {
			for _, t := range list {
				tasks[t.ID] = t
			}
		}
	}
}

func saveData() {
	mutex.Lock()
	defer mutex.Unlock()
	// Projects
	var plist []Project
	for _, p := range projects {
		plist = append(plist, p)
	}
	pdata, _ := json.MarshalIndent(plist, "", "  ")
	_ = ioutil.WriteFile(projectsFile, pdata, 0644)
	// Tasks
	var tlist []Task
	for _, t := range tasks {
		tlist = append(tlist, t)
	}
	data, _ := json.MarshalIndent(tlist, "", "  ")
	_ = ioutil.WriteFile(tasksFile, data, 0644)
} 