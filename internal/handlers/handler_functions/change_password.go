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
			userid, data_bool := check_if_authorize(r)
			var requestData struct {
				Email       string `json:"email"`
				NewPassword string `json:"password"`
				OldPassword string `json:"old_password,omitempty"`
			}

			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				log.Printf("Error decode JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if requestData.NewPassword == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ip := r.RemoteAddr
			if requestData.OldPassword != "" {
				if data_bool {
					email, err := queries.ValidateOldPassword(pool.PoolConns, userid, requestData.OldPassword)
					log.Println(email)
					if err != nil || email == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					err = queries.ChangePasswordQuery(pool.PoolConns, email, requestData.NewPassword)
				} else {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				if requestData.Email == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// Логика для случая без старого пароля
				if !canChange[requestData.Email] {
					w.WriteHeader(http.StatusConflict)
					return
				}
				if codeCache[requestData.Email].expiresAt.Before(time.Now()) {
					canChange[requestData.Email] = false
					w.WriteHeader(http.StatusLocked)
					return
				}
				if attemptCache[ip].lastTry.After(time.Now()) {
					canChange[requestData.Email] = false
					w.WriteHeader(http.StatusRequestTimeout)
					return
				}
				err = queries.ChangePasswordQuery(pool.PoolConns, requestData.Email, requestData.NewPassword)
			}

			if err != nil {
				log.Printf("Error change password in db: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			canChange[requestData.Email] = false
			codeCache[requestData.Email] = codeInfo{
				code:      codeCache[requestData.Email].code,
				expiresAt: time.Time{},
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	return http.HandlerFunc(fn)
}
