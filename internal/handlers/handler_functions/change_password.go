package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func ChangePassword(pool *pool_conections.Pool_conections) http.Handler {
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
				w.WriteHeader(http.StatusBadRequest)
				return
			}


			ip := r.RemoteAddr
			if !canChange[creds.Email] {
				w.WriteHeader(http.StatusConflict)
				return
			}
			if codeCache[creds.Email].expiresAt.Before(time.Now()) {
				canChange[creds.Email] = false
				w.WriteHeader(http.StatusLocked)
				return
			}
			if attemptCache[ip].lastTry.After(time.Now()) {
				canChange[creds.Email] = false
				w.WriteHeader(http.StatusRequestTimeout)
				return
			}

			err = queries.ChangePasswordQuery(pool.PoolConns, creds.Email, creds.NewPassword)
			if err != nil {
				log.Printf("Error change password in db: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			canChange[creds.Email] = false
			codeCache[creds.Email] = codeInfo{
				code:      codeCache[creds.Email].code,
				expiresAt: time.Time{},
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	return http.HandlerFunc(fn)
}
