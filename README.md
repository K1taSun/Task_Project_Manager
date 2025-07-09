# 🚀 Task Project Manager

**Nowoczesny, kompletny menadżer projektów i zadań** z backendem w Go oraz responsywnym frontendem HTML+JS. Pozwala wygodnie zarządzać projektami, zadaniami, filtrować, edytować, eksportować dane i korzystać z motywu jasnego/ciemnego.

---

## ✨ Funkcje

- **Zarządzanie projektami**: dodawanie, edycja, usuwanie
- **Zarządzanie zadaniami**: tytuł, deadline, tagi, priorytet (0-5), status wykonania
- **Filtrowanie i sortowanie** zadań (po tagu, priorytecie, dacie)
- **Eksport danych** do JSON i CSV jednym kliknięciem
- **Walidacja danych** po stronie backendu
- **Motyw jasny/ciemny** (przełącznik)
- **Responsywny, nowoczesny interfejs** z Material Icons
- **Trwałość danych** w plikach JSON (brak bazy danych)

---

## 🖥️ Technologie

- **Backend:** Go 1.24 (REST API, walidacja, eksport, persystencja do plików JSON)
- **Frontend:** HTML5, CSS3 (custom design), JavaScript (fetch API, dynamiczne renderowanie)
- **UI:** Material Icons, motyw jasny/ciemny, responsywność

---

## 📦 Struktura projektu

```
├── main.go           # Start serwera, routing
├── handlers.go       # Endpointy REST API
├── models.go         # Struktury danych
├── storage.go        # Wczytywanie/zapisywanie danych
├── validation.go     # Walidacja danych
├── index.html        # Frontend (HTML, CSS, JS)
├── data_projects.json# Projekty (dane)
├── data_tasks.json   # Zadania (dane)
```

---

## 🔗 API (REST)

- `GET    /projects` — lista projektów
- `POST   /projects` — dodaj projekt `{ name }`
- `GET    /projects/{id}` — szczegóły projektu
- `PUT    /projects/{id}` — edycja projektu
- `DELETE /projects/{id}` — usuń projekt
- `GET    /tasks` — lista zadań (filtrowanie: `tag`, `min_priority`, `max_priority`, `before`, `after`, `sort`, `order`)
- `POST   /tasks` — dodaj zadanie
- `GET    /tasks/{id}` — szczegóły zadania
- `PUT    /tasks/{id}` — edycja zadania
- `DELETE /tasks/{id}` — usuń zadanie
- `GET    /projects/{id}/tasks` — zadania w projekcie
- `POST   /projects/{id}/tasks` — dodaj zadanie do projektu
- `GET    /export?format=json|csv` — eksport wszystkich danych

**Przykład zadania:**
```json
{
  "id": 1,
  "project_id": 123,
  "title": "Zaimplementuj eksport",
  "deadline": "2025-07-31T00:00:00Z",
  "tags": ["backend", "eksport"],
  "priority": 3,
  "done": false
}
```

---

## 🏁 Uruchomienie

1. **Wymagania:** Go 1.24+
2. **Start backendu:**
   ```bash
   go run main.go
   ```
   Serwer ruszy na `localhost:8080`
3. **Frontend:** Otwórz `index.html` w przeglądarce

---

## 📂 Przykładowe dane

**data_projects.json**
```json
[
  { "id": 428131, "name": "Buffer_MacOS" },
  { "id": 662013, "name": "MacOS" }
]
```

**data_tasks.json**
```json
[
  {
    "id": 193891,
    "project_id": 0,
    "title": "XD",
    "deadline": "2025-07-31T00:00:00Z",
    "tags": ["11"],
    "priority": 3,
    "done": false
  }
]
```

---

## ⚖️ Licencja

Projekt na licencji MIT. Możesz używać, modyfikować i rozwijać bez ograniczeń.

---

> **Autor:** Nikita Parkovskyi

