package main

import (
	"log"
	"net/http"
	"os"
)

// main.go - główny plik programu
// tutaj startuje cała aplikacja

func main() {
	// ustawiamy logi
	log.Println("Startuje Task Project Manager...")

	// wczytujemy konfigurację
	config := loadConfigFromEnv()
	printConfig(config)

	// sprawdzamy czy konfiguracja jest poprawna
	err := validateConfig(config)
	if err != nil {
		log.Fatal("Blad konfiguracji:", err)
	}

	// sprawdzamy czy pliki istnieją
	// jak nie to tworzymy
	if _, err := os.Stat(config.ProjectsFile); os.IsNotExist(err) {
		log.Println("Tworze plik", config.ProjectsFile)
		err := SaveProjects()
		if err != nil {
			log.Fatal("Blad przy tworzeniu pliku projektow:", err)
		}
	}

	if _, err := os.Stat(config.TasksFile); os.IsNotExist(err) {
		log.Println("Tworze plik", config.TasksFile)
		err := SaveTasks()
		if err != nil {
			log.Fatal("Blad przy tworzeniu pliku zadan:", err)
		}
	}

	// wczytujemy dane
	err = LoadProjects()
	if err != nil {
		log.Fatal("Blad wczytywania projektow:", err)
	}
	err = LoadTasks()
	if err != nil {
		log.Fatal("Blad wczytywania zadan:", err)
	}

	// dodajemy endpointy
	http.HandleFunc("/projects", withCORS(logMiddleware(projectsHandler)))
	http.HandleFunc("/projects/", withCORS(logMiddleware(projectHandler)))
	http.HandleFunc("/tasks", withCORS(logMiddleware(tasksHandler)))
	http.HandleFunc("/tasks/", withCORS(logMiddleware(taskHandler)))
	http.HandleFunc("/export", withCORS(logMiddleware(exportHandler)))

	// nowe endpointy
	http.HandleFunc("/badges", withCORS(logMiddleware(badgesHandler)))
	http.HandleFunc("/stats", withCORS(logMiddleware(statsHandler)))
	http.HandleFunc("/ws", logMiddleware(wsHandler))

	// strona glowna
	http.HandleFunc("/", logMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
		} else {
			http.NotFound(w, r)
		}
	}))

	// uruchamiamy serwer
	serverAddr := config.Host + ":" + config.Port
	log.Println("Serwer dziala na", serverAddr)
	log.Println("Otworz http://" + serverAddr + " w przegladarce")

	// startujemy serwer
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Blad uruchamiania serwera:", err)
	}
}
