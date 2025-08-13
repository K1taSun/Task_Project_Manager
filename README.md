# Task Project Manager

**Prosty menadżer projektów i zadań** napisany w Go. Aplikacja pozwala zarządzać projektami i zadaniami.

---

## Funkcje

- **Projekty**: dodawanie, edycja, usuwanie
- **Zadania**: tytuł, deadline, tagi, priorytet (0-5), status wykonania
- **Filtrowanie** zadań (po tagu, priorytecie, dacie)
- **Eksport** do JSON i CSV
- **Walidacja** danych
- **Motyw jasny/ciemny**
- **Responsywny interfejs**

---

## Technologie

- **Backend:** Go 1.24
- **Frontend:** HTML, CSS, JavaScript
- **Dane:** pliki JSON

---

## Struktura projektu

```
├── main.go           # Start serwera
├── handlers.go       # API endpoints
├── models.go         # Struktury danych
├── storage.go        # Zapis/odczyt danych
├── validation.go     # Walidacja
├── utils.go          # Funkcje pomocnicze
├── index.html        # Frontend
├── data_projects.json# Projekty
├── data_tasks.json   # Zadania
```

---

## API

- `GET    /projects` — lista projektów
- `POST   /projects` — dodaj projekt
- `GET    /projects/{id}` — szczegóły projektu
- `PUT    /projects/{id}` — edycja projektu
- `DELETE /projects/{id}` — usuń projekt
- `GET    /tasks` — lista zadań
- `POST   /tasks` — dodaj zadanie
- `GET    /tasks/{id}` — szczegóły zadania
- `PUT    /tasks/{id}` — edycja zadania
- `DELETE /tasks/{id}` — usuń zadanie
- `GET    /export?format=json|csv` — eksport danych

---

## Uruchomienie

1. **Wymagania:** Go 1.24+
2. **Start:**
   ```bash
   go run main.go
   ```
3. **Otwórz:** `http://localhost:8080`

---

## Przykładowe dane

**data_projects.json**
```json
[
  { "id": 1, "name": "Projekt 1" },
  { "id": 2, "name": "Projekt 2" }
]
```

**data_tasks.json**
```json
[
  {
    "id": 1,
    "project_id": 1,
    "title": "Zadanie 1",
    "deadline": "2025-07-31T00:00:00Z",
    "tags": ["tag1"],
    "priority": 3,
    "done": false
  }
]
```

---

## TODO

- [ ] Dodać bazę danych
- [ ] Dodać użytkowników
- [ ] Dodać więcej walidacji
- [ ] Poprawić błędy

---

## Licencja

MIT

---

**Autor:** Student 3 roku

