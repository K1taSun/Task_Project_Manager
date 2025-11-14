package app

import (
	"encoding/json"
	"log"
	"os"
)

func LoadProjects() error {
	projectsFile := GetCurrentConfig().ProjectsFile
	file, err := os.Open(projectsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik projektow nie istnieje, tworze nowy")
			return SaveProjects()
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
	for _, p := range list {
		projects[p.ID] = p
		if p.ID > maxID {
			maxID = p.ID
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
	return os.WriteFile(GetCurrentConfig().ProjectsFile, data, 0o644)
}

func LoadTasks() error {
	tasksFile := GetCurrentConfig().TasksFile
	file, err := os.Open(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik zadan nie istnieje, tworze nowy")
			return SaveTasks()
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
	for _, t := range list {
		tasks[t.ID] = t
		if t.ID > maxID {
			maxID = t.ID
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
	return os.WriteFile(GetCurrentConfig().TasksFile, data, 0o644)
}

func LoadUsers() error {
	usersFile := GetCurrentConfig().UsersFile
	file, err := os.Open(usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Plik użytkowników nie istnieje, tworzę nowy")
			return SaveUsers()
		}
		return err
	}
	defer file.Close()

	var list []User
	if err := json.NewDecoder(file).Decode(&list); err != nil {
		return err
	}

	userMutex.Lock()
	defer userMutex.Unlock()
	users = make(map[int]User)
	maxID := 0
	for _, u := range list {
		users[u.ID] = u
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	nextUserID = maxID + 1
	return nil
}

func SaveUsers() error {
	userMutex.RLock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	userMutex.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GetCurrentConfig().UsersFile, data, 0o644)
}
