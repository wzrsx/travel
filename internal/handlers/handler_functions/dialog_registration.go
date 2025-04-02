package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

var (
	pendingRegistrations = make(map[string]pendingRegistration)
	registrationMutex    sync.Mutex
)

type pendingRegistration struct {
	Username string
	Password string
	Email    string
	Pool     *pool_conections.Pool_conections
}

func Registration(p *pool_conections.Pool_conections) http.Handler {
	registration := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "Invalid Content-Type: ", http.StatusUnsupportedMediaType)
			return
		}

		// Декодируем JSON из тела запроса
		type CredentialsRegistration struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var creds CredentialsRegistration
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if creds.Email == "" || creds.Password == "" || creds.Username == "" {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Проверка наличия почты в базе данных
		err = queries.ExistsEmail(creds.Email, p.PoolConns)
		if err != nil {
			if err.Error() == "email exists" {
				http.Error(w, "Email exists", http.StatusConflict)
			}
		}

		registrationMutex.Lock()
		pendingRegistrations[creds.Email] = pendingRegistration{
			Username: creds.Username,
			Password: creds.Password,
			Email:    creds.Email,
			Pool:     p,
		}
		registrationMutex.Unlock()

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(registration)
}
