package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func CheckEmailIntoDB(pool *pool_conections.Pool_conections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Декодируем JSON из тела запроса
			type CredentialsCheckEmail struct {
				Email string `json:"email"`
			}

			var creds CredentialsCheckEmail
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
				return
			}
			err = queries.ExistsEmail(creds.Email, pool.PoolConns)
			switch err.Error() {
			case "email exists":
				w.WriteHeader(http.StatusOK)
				return
			case "OK":
				w.WriteHeader(409)
				return
			default:
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("error checking exist email: %v", err)
				http.Error(w, "error checking exist email", http.StatusBadRequest)
			}
		}
	}
	return http.HandlerFunc(fn)
}
