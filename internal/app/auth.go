package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "session_token"
	sessionDuration   = 24 * time.Hour
)

var (
	sessionStore = make(map[string]session)
	sessionMutex sync.RWMutex

	roleHierarchy = map[string]int{
		"viewer":  1,
		"manager": 2,
		"admin":   3,
	}
)

type session struct {
	UserID    int
	ExpiresAt time.Time
}

type contextKey string

const userContextKey contextKey = "authUser"

// EnsureDefaultAdmin tworzy użytkownika administratora, jeśli żaden nie istnieje.
func EnsureDefaultAdmin() error {
	if adminCount() > 0 {
		return nil
	}

	password := "admin123"
	_, err := CreateUser("admin@example.com", "Administrator", password, "admin")
	if err != nil {
		return err
	}
	log.Println("Utworzono domyślnego użytkownika admin:")
	log.Println("  Email: admin@example.com")
	log.Println("  Hasło: admin123")
	log.Println("Zmień dane logowania po pierwszym uruchomieniu.")
	return nil
}

// CreateUser dodaje nowego użytkownika z zahashowanym hasłem.
func CreateUser(email, name, password, role string) (User, error) {
	email = normalizeEmail(email)
	role = strings.ToLower(strings.TrimSpace(role))
	if email == "" || password == "" {
		return User{}, errors.New("email i hasło są wymagane")
	}
	if !isValidRole(role) {
		return User{}, errors.New("niepoprawna rola użytkownika")
	}

	userMutex.Lock()
	if existing, ok := getUserByEmailLocked(email); ok {
		userMutex.Unlock()
		return User{}, errors.New("użytkownik o podanym emailu już istnieje: " + existing.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		userMutex.Unlock()
		return User{}, err
	}

	now := time.Now()
	user := User{
		ID:           generateUserID(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	users[user.ID] = user
	userMutex.Unlock()

	if err := SaveUsers(); err != nil {
		return User{}, err
	}

	return user, nil
}

// AuthenticateUser weryfikuje dane logowania.
func AuthenticateUser(email, password string) (*User, error) {
	email = normalizeEmail(email)
	userMutex.RLock()
	user, ok := getUserByEmailLocked(email)
	userMutex.RUnlock()

	if !ok {
		return nil, errors.New("niepoprawny email lub hasło")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("niepoprawny email lub hasło")
	}

	return &user, nil
}

// GetUserByID zwraca użytkownika po ID.
func GetUserByID(id int) (User, bool) {
	userMutex.RLock()
	defer userMutex.RUnlock()
	u, ok := users[id]
	return u, ok
}

// GetAllUsers zwraca listę użytkowników (bez hashy).
func GetAllUsers() []User {
	return sanitizeUsers(getAllUsers())
}

func getAllUsers() []User {
	userMutex.RLock()
	defer userMutex.RUnlock()
	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, u)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func getUserByEmailLocked(email string) (User, bool) {
	for _, u := range users {
		if u.Email == email {
			return u, true
		}
	}
	return User{}, false
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func isValidRole(role string) bool {
	_, ok := roleHierarchy[role]
	return ok
}

func sanitizeUser(u User) User {
	u.PasswordHash = ""
	return u
}

func sanitizeUsers(list []User) []User {
	safe := make([]User, len(list))
	for i, u := range list {
		safe[i] = sanitizeUser(u)
	}
	return safe
}

func hasAnotherAdmin(excludeID int) bool {
	userMutex.RLock()
	defer userMutex.RUnlock()
	for _, u := range users {
		if u.Role == "admin" && u.ID != excludeID {
			return true
		}
	}
	return false
}

func adminCount() int {
	userMutex.RLock()
	defer userMutex.RUnlock()
	count := 0
	for _, u := range users {
		if u.Role == "admin" {
			count++
		}
	}
	return count
}

// Sesje

func createSession(userID int) (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	sessionMutex.Lock()
	sessionStore[token] = session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	sessionMutex.Unlock()
	return token, nil
}

func getSession(token string) (session, bool) {
	sessionMutex.RLock()
	sess, ok := sessionStore[token]
	sessionMutex.RUnlock()
	if !ok {
		return session{}, false
	}
	if sess.ExpiresAt.Before(time.Now()) {
		deleteSession(token)
		return session{}, false
	}
	return sess, true
}

func deleteSession(token string) {
	sessionMutex.Lock()
	delete(sessionStore, token)
	sessionMutex.Unlock()
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// Middleware

func authenticateRequest(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, errors.New("brak sesji")
	}

	sess, ok := getSession(cookie.Value)
	if !ok {
		return nil, errors.New("sesja wygasła")
	}

	user, exists := GetUserByID(sess.UserID)
	if !exists {
		return nil, errors.New("użytkownik nie istnieje")
	}
	return &user, nil
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		user, err := authenticateRequest(r)
		if err != nil {
			clearSessionCookie(w)
			writeJSONError(w, http.StatusUnauthorized, "Wymagane logowanie")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func requireRole(minRole string, next http.HandlerFunc) http.HandlerFunc {
	return requireAuth(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user == nil {
			writeJSONError(w, http.StatusUnauthorized, "Wymagane logowanie")
			return
		}

		if !hasRoleOrHigher(user.Role, minRole) {
			writeJSONError(w, http.StatusForbidden, "Brak uprawnień")
			return
		}
		next(w, r)
	})
}

func hasRoleOrHigher(role, minRole string) bool {
	current, okCurrent := roleHierarchy[role]
	required, okRequired := roleHierarchy[minRole]
	if !okCurrent || !okRequired {
		return false
	}
	return current >= required
}

func GetUserFromContext(r *http.Request) *User {
	user, _ := r.Context().Value(userContextKey).(*User)
	return user
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionDuration.Seconds()),
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func deleteSessionsForUser(userID int) {
	sessionMutex.Lock()
	for token, sess := range sessionStore {
		if sess.UserID == userID {
			delete(sessionStore, token)
		}
	}
	sessionMutex.Unlock()
}

// VerifyAdminRegistrationToken sprawdza, czy podany token jest poprawny dla rejestracji admina.
func VerifyAdminRegistrationToken(token string) bool {
	cfg := GetCurrentConfig()
	if cfg.AdminRegistrationToken == "" {
		return false
	}
	return token == cfg.AdminRegistrationToken
}
