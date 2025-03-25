package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"travel/internal/handlers/JWT"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

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

		// Делаем запрос в БД
		user_reg_res := queries.NewUserRegistrationResult()
		err = user_reg_res.RegistrationQuery(creds.Username, creds.Email, creds.Password, p.PoolConns)
		if err != nil {
			if err.Error() == "email exists" {
				http.Error(w, "Email exists", http.StatusConflict)
			}
			log.Printf("Registration Querry returns error: %v", err)
			http.Error(w, "Registration Querry error", http.StatusBadRequest)
			return
		}

		// устанавливаем Cookie с JWT
		err = JWT.SetCookeJWT(w, user_reg_res.UserID, user_reg_res.Username)
		if err != nil {
			http.Error(w, "Error set cookie JWT", http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(registration)
}