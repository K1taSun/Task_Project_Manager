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
- **Zaawansowane metryki** – dashboard z podsumowaniami, wykresami priorytetów, przepływu pracy i obciążenia wykonawców
- **Autoryzacja i role** – logowanie z użyciem ciasteczek HttpOnly, role `admin`/`manager`/`viewer`, zarządzanie użytkownikami

---

## Interfejs w skrócie

- Górny pasek zawiera przycisk **Szybkiego podglądu** (menu), eksport danych oraz przełącznik motywu.
- Lewy panel pokazuje projekty z liczbą zadań i podpowiedzią edycji pod prawym klawiszem.
- W widoku projektu od razu widać listę zadań z akcjami „oznacz ukończone” i „edytuj”.
- Filtry (wyszukiwarka, tag, priorytet, dzień) działają na bieżącej liście zadań.
- Przycisk typu FAB w prawym dolnym rogu umożliwia dodanie zadania z dowolnego miejsca.

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
├── data_projects.json     # bieżące dane projektów (generowane przy starcie)
├── data_tasks.json        # bieżące dane zadań (generowane przy starcie)
├── data_users.json        # bieżące dane użytkowników (tworzone automatycznie)
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
5. **Logowanie:** przy pierwszym uruchomieniu tworzony jest użytkownik `admin@example.com` z hasłem `admin123`. Zmień je od razu po zalogowaniu.
6. **Dashboard metryk:** kliknij ikonę wykresu w górnym pasku, aby zobaczyć podsumowania i wykresy (Chart.js jest ładowany z CDN).

> Pliki `data_projects.json`, `data_tasks.json` oraz `data_users.json` są wykluczone z repozytorium (zapisują bieżący stan). Jeśli zaczynasz od pustej instalacji, aplikacja sama utworzy te pliki przy pierwszym uruchomieniu.

---

## Uwierzytelnianie i role

- Sesje trzymane są w pamięci serwera i identyfikowane ciasteczkiem `session_token` (HttpOnly, Lax).
- API:
  - `POST /auth/login` – logowanie (JSON: `email`, `password`)
  - `POST /auth/logout` – wylogowanie
  - `GET /auth/me` – informacja o aktualnym użytkowniku
  - `GET /users` – lista użytkowników (tylko `admin`)
  - `POST /users` – tworzenie użytkownika (tylko `admin`)
  - `PUT /users/{id}` / `DELETE /users/{id}` – zmiana/usuń (tylko `admin`)
- Role:
  - `admin` – pełny dostęp do danych i zarządzania kontami
  - `manager` – zarządzanie projektami i zadaniami
  - `viewer` – dostęp tylko do odczytu (UI ukrywa akcje modyfikujące)

Frontend wykrywa rolę po zalogowaniu i ukrywa akcje wymagające uprawnień; próba wykonania niedozwolonej operacji zakończy się błędem 403 z backendu.

---

## Dashboard metryk

Panel dostępny z poziomu górnego paska (`ikona analytics`) prezentuje:

- **Podsumowanie** – stopień ukończenia, zaległe terminy, średni priorytet, liczbę ukończonych zadań w ostatnich 7 dniach.
- **Dystrybucję priorytetów** – wykres kołowy w oparciu o bieżące dane.
- **Tempo pracy (14 dni)** – wykres liniowy pokazujący utworzone i ukończone zadania.
- **Obciążenie wykonawców** – porównanie otwartych i zamkniętych zadań per osoby.
- **Najczęstsze tagi** – szybki podgląd dominujących etykiet w projektach.

Dane są odświeżane automatycznie przy zmianach (REST + WebSocket). Tryb ciemny/jasny aktualizuje styl wykresów bez przeładowania.

---

## Konfiguracja

| Zmienna          | Domyślna wartość      | Opis                                                |
|------------------|------------------------|-----------------------------------------------------|
| `HOST`           | `localhost`            | Adres interfejsu serwera HTTP                      |
| `PORT`           | `8080`                 | Port nasłuchu HTTP                                 |
| `PROJECTS_FILE`  | `data_projects.json`   | Ścieżka do pliku z projektami                      |
| `TASKS_FILE`     | `data_tasks.json`      | Ścieżka do pliku z zadaniami                       |
| `USERS_FILE`     | `data_users.json`      | Ścieżka do pliku z użytkownikami                   |
| `SESSION_SECRET` | losowo generowany      | Sekret wykorzystywany do zabezpieczenia sesji      |

Pliki z danymi tworzą się automatycznie. Jeśli chcesz zacząć od pustej bazy, ustaw w nich `[]` lub wskaż inne lokalizacje.

---

## Przykładowe dane

Pliki `data_projects.json` i `data_tasks.json` są generowane automatycznie przy pierwszym uruchomieniu. Poniżej przykładowa struktura rekordów:

```jsonc
// data_projects.json
[
  {
    "id": 1,
    "name": "Projekt 1",
    "status": "active",
    "badge": { "text": "Nowy", "background": "#3b82f6", "type": "status" },
    "created_at": "2025-11-02T20:15:14+01:00",
    "updated_at": "2025-11-02T20:15:14+01:00"
  }
]

// data_tasks.json
[
  {
    "id": 1,
    "project_id": 1,
    "title": "Przygotować plan sprintu",
    "priority": 4,
    "tags": ["planowanie"],
    "deadline": "2025-11-10T00:00:00Z",
    "badge": { "text": "Krytyczny", "background": "#dc2626", "type": "priority" },
    "created_at": "2025-11-02T20:16:46+01:00",
    "updated_at": "2025-11-02T20:16:46+01:00"
  }
]
```

> Tip: Przed publikacją w repozytorium możesz zresetować pliki danych do pustych tablic, aby nie trzymać prywatnych informacji.

---

## TODO / pomysły na rozwój

- [ ] Integracja z relacyjną bazą danych (PostgreSQL/SQLite)
- [x] Autoryzacja użytkowników i role
- [x] Zaawansowane metryki oraz dashboard
- [ ] Powiadomienia e-mail / webhooki
- [ ] Testy jednostkowe i automatyczna walidacja linterem/go vet

---

## Licencja

MIT

---

**Autor:** Nikita Parkovskyi

