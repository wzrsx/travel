package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"travel/internal/handlers/JWT"
	"travel/internal/repositories/queries"
)

func CheckPassCode() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type CredentialsSendEmail struct {
				Email string `json:"email"`
				Code  string `json:"code"`
				IsReg bool   `json:"isreg"`
			}
			ip := r.RemoteAddr

			var creds CredentialsSendEmail
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				log.Printf("Failed to decode JSON: %v", err)
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
			if creds.IsReg {
				registrationMutex.Lock()
				pending, exists := pendingRegistrations[creds.Email]
				registrationMutex.Unlock()

				if !exists {
					http.Error(w, "No pending registration for this email", http.StatusNotFound)
					return
				}

				user_reg_res := queries.NewUserRegistrationResult()
				err = user_reg_res.RegistrationQuery(pending.Username, pending.Email, pending.Password, pending.Pool.PoolConns)
				if err != nil {
					if err.Error() == "email exists" {
						http.Error(w, "Email exists", http.StatusConflict)
					}
					log.Printf("Registration Querry returns error: %v", err)
					http.Error(w, "Registration Querry error", http.StatusBadRequest)
					return
				}

				// Устанавливаем куку
				err = JWT.SetCookeJWT(w, user_reg_res.UserID, user_reg_res.Username)
				if err != nil {
					http.Error(w, "Error set cookie JWT", http.StatusUnauthorized)
				}
			}
			canChange[creds.Email] = true
			codeCache[creds.Email] = codeInfo{
				code:      codeCache[creds.Email].code,
				expiresAt: time.Now().Add(time.Minute * 10),
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	return http.HandlerFunc(fn)
}
