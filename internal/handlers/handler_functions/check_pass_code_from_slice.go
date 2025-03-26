package handler_functions

import (
	"encoding/json"
	"net/http"
	"time"
)

func CheckPassCode() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type CredentialsSendEmail struct {
				Email string `json:"email"`
				Code  string `json:"code"`
			}
			ip := r.RemoteAddr

			var creds CredentialsSendEmail
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
				return
			}

			if attemptCache[ip].attempts >= 10 {
				attemptCache[ip] = attemptInfo{
					attempts: 0,
					lastTry:  attemptCache[ip].lastTry.Add(time.Minute * 1),
				}
			}
			if attemptCache[ip].lastTry.After(time.Now()) {
				unlockTime := attemptCache[ip].lastTry
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusLocked)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"unlock_at": unlockTime.Format(time.RFC3339),
				})
				return
			}
			if creds.Code != codeCache[creds.Email].code {
				attemptCache[ip] = attemptInfo{
					attempts: attemptCache[ip].attempts + 1,
					lastTry:  time.Now(),
				}
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	return http.HandlerFunc(fn)
}
