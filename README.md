# Task Project Manager

**Task Project Manager** to lekka aplikacja webowa pozwalająca zarządzać projektami i zadaniami. Backend napisany w Go udostępnia REST API i prosty mechanizm WebSocket, a frontend w czystym HTML/JS zapewnia nowoczesny interfejs z szybkim podglądem i filtrowaniem.

---

## Funkcje

- **Projekty** – tworzenie, edycja, archiwizacja/usuwanie, odznaki statusowe
- **Zadania** – tytuł, opis, deadline, tagi, priorytet (0-5), czas szacowany/rzeczywisty, przypisanie do projektu (również przy edycji)
- **Panel szybkiego podglądu** – wyszukiwarka wszystkich projektów i zadań dostępna z górnego paska
- **Widok projektu** – natychmiastowa lista zadań danego projektu, szybkie akcje (oznacz ukończone, edycja, dodaj pierwsze zadanie)
- **Filtrowanie i sortowanie** – m.in. po tagu, priorytecie, dacie
- **Eksport** danych do JSON i CSV jednym kliknięciem
- **Walidacja** danych wejściowych po stronie backendu
- **WebSocket** – push powiadomień o zmianach do otwartych kart
- **Motyw jasny/ciemny** oraz responsywny interfejs przygotowany na urządzenia mobilne

---

## Technologie

- **Backend:** Go 1.21+
- **Frontend:** HTML, CSS, JavaScript
- **Dane:** pliki JSON

---

## Struktura projektu

```
├── cmd/
│   └── task-manager/
│       └── main.go        # punkt startowy aplikacji
├── internal/
│   └── app/
│       ├── config.go      # konfiguracja, zmienne środowiskowe
│       ├── handlers.go    # REST API, WebSocket, middleware
│       ├── models.go      # struktury danych + pamięciowe „repozytorium”
│       ├── storage.go     # zapisy/odczyty JSON
│       ├── utils.go       # narzędzia, logger, odpowiedzi JSON
│       └── validation.go  # walidacja wpisów
├── index.html             # jedyny plik frontendu
├── data_projects.json     # plik z projektami
├── data_tasks.json        # plik z zadaniami
├── go.mod / go.sum        # moduł Go
└── README.md
```

---

- **Projekty**
  - `GET    /projects` – lista projektów
  - `POST   /projects` – nowy projekt
  - `GET    /projects/{id}` – szczegóły projektu
  - `PUT    /projects/{id}` – edycja projektu
  - `DELETE /projects/{id}` – usunięcie projektu
  - `GET    /projects/{id}/tasks` – zadania projektu, `POST` – dodanie zadania do projektu
- **Zadania**
  - `GET    /tasks` – lista zadań z filtrowaniem (tag, priorytet, zakres dat, sortowanie)
  - `POST   /tasks` – utworzenie zadania (opcjonalnie bez projektu)
  - `GET    /tasks/{id}` – podgląd zadania
  - `PUT    /tasks/{id}` – aktualizacja, w tym przeniesienie do innego projektu
  - `DELETE /tasks/{id}` – usunięcie
- **Inne**
  - `GET    /badges` – słownik odznak dla frontendu
  - `GET    /export?format=json|csv` – eksport wszystkich danych
  - `GET    /ws` – WebSocket z powiadomieniami o zmianach

---

## Uruchomienie

1. **Wymagania:** Go 1.21+, przeglądarka z włączonym JavaScript.
2. **Instalacja zależności:**
   ```bash
   go mod tidy
   ```
3. **Start aplikacji:**
   ```bash
   go run ./cmd/task-manager
   ```
   Domyślnie backend nasłuchuje na `localhost:8080`. Zmienne środowiskowe (`HOST`, `PORT`, `PROJECTS_FILE`, `TASKS_FILE`) pozwalają nadpisać konfigurację.
4. **Frontend:** otwórz w przeglądarce `http://localhost:8080`.

---

## Przykładowe dane

Pliki `data_projects.json` i `data_tasks.json` są generowane automatycznie przy pierwszym uruchomieniu. Możesz je wyczyścić do pustej tablicy `[]`, jeżeli chcesz zacząć od zera.

---

## TODO / pomysły na rozwój

- [ ] Integracja z relacyjną bazą danych (PostgreSQL/SQLite)
- [ ] Autoryzacja użytkowników i role
- [ ] Zaawansowane metryki oraz dashboard
- [ ] Powiadomienia e-mail / webhooki

---

## Licencja

MIT

---

**Autor:** Student 3 roku

