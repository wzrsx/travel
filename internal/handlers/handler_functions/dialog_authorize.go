package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"travel/internal/handlers/JWT"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func Authorize(p *pool_conections.Pool_conections) http.Handler {
	authorize := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем Content-Type
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "Invalid Content-Type: ", http.StatusUnsupportedMediaType)
			return
		}

		// Декодируем JSON из тела запроса
		type CredentialsAuthorize struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var creds CredentialsAuthorize
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if creds.Email == "" || creds.Password == "" {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Делаем запрос в БД
		user_auth_res := queries.NewUserAuthResult()
		err = user_auth_res.AuthorizeQuery(creds.Email, creds.Password, p.PoolConns)
		if err != nil {
			log.Printf("Authorize Querry returns error: %v", err)
			if err.Error() == "email not found" {
				http.Error(w, "Invalid email", http.StatusNotFound)
			}else if err.Error() == "pass invalid"{
				http.Error(w, "Invalid email", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// устанавливаем Cookie с JWT
		err = JWT.SetCookeJWT(w, user_auth_res.UserID, user_auth_res.Username)
		if err != nil {
			http.Error(w, "Error set cookie JWT", http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(authorize)
}