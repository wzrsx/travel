package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func ChangePassword() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type Credentials struct {
				Email       string `json:"email"`
				NewPassword string `json:"password"`
			}
			var creds Credentials

			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				log.Printf("Error decode JSON: %v", err)
				return
			}

			ip := r.RemoteAddr
			if !canChange[creds.Email] {
				w.WriteHeader(http.StatusConflict)
				return
			}
			if codeCache[creds.Email].expiresAt.Before(time.Now()) {
				w.WriteHeader(http.StatusLocked)
			}
			if attemptCache[ip].lastTry.After(time.Now()) {
				w.WriteHeader(http.StatusRequestTimeout)
			}

			// Запрос на изменение пароля
			canChange[creds.Email] = false
			w.WriteHeader(http.StatusOK)
		}
	}
	return http.HandlerFunc(fn)
}
