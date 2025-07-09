package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var wsClients = make(map[*websocket.Conn]bool)
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Pozwól na połączenia z dowolnego origin
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Dodaj klienta do mapy
	wsClients[conn] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(wsClients))

	defer func() {
		delete(wsClients, conn)
		conn.Close()
		log.Printf("WebSocket client disconnected. Total clients: %d", len(wsClients))
	}()

	// Nasłuchuj na wiadomości
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func broadcastChange() {
	if len(wsClients) == 0 {
		return
	}

	message := []byte("update")
	for client := range wsClients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			client.Close()
			delete(wsClients, client)
		}
	}
}

func main() {
	// Konfiguracja logowania
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Task Project Manager...")

	// Sprawdź czy pliki danych istnieją, jeśli nie - utwórz je
	if _, err := os.Stat("data_projects.json"); os.IsNotExist(err) {
		log.Println("Creating data_projects.json")
		if err := SaveProjects(); err != nil {
			log.Fatalf("Failed to create projects file: %v", err)
		}
	}
	if _, err := os.Stat("data_tasks.json"); os.IsNotExist(err) {
		log.Println("Creating data_tasks.json")
		if err := SaveTasks(); err != nil {
			log.Fatalf("Failed to create tasks file: %v", err)
		}
	}

	if err := LoadProjects(); err != nil {
		log.Fatalf("Błąd wczytywania projektów: %v", err)
	}
	if err := LoadTasks(); err != nil {
		log.Fatalf("Błąd wczytywania zadań: %v", err)
	}

	// Rejestracja endpointów z middleware
	http.HandleFunc("/projects", withCORS(logMiddleware(projectsHandler)))
	http.HandleFunc("/projects/", withCORS(logMiddleware(projectHandler)))
	http.HandleFunc("/tasks", withCORS(logMiddleware(tasksHandler)))
	http.HandleFunc("/tasks/", withCORS(logMiddleware(taskHandler)))
	http.HandleFunc("/export", withCORS(logMiddleware(exportHandler)))
	http.HandleFunc("/ws", logMiddleware(wsHandler))

	// Serwuj pliki statyczne
	http.HandleFunc("/", logMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
		} else {
			http.NotFound(w, r)
		}
	}))

	log.Println("Server running on :8080")
	log.Println("Open http://localhost:8080 in your browser")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
