package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func NewReviewHandler(p *pool_conections.Pool_conections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/new_review.html")
		if err != nil {
			log.Printf("%s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if data.UserID == 0 {
			http.Redirect(w, r, "/?openLoginDialog=true", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodGet {
			// Проверяем наличие route_id в URL
			routeID := r.URL.Query().Get("route_id")
			if routeID == "" {
				// Если route_id отсутствует, перенаправляем на главную страницу
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				log.Printf("Failed to render template: %s", err.Error())
			}
			return
		}

		if r.Method == http.MethodPost {
			// Проверяем Content-Type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Invalid Content-Type: expected application/json", http.StatusUnsupportedMediaType)
				return
			}

			// Декодируем JSON
			type CredentialsNewReview struct {
				Username    string `json:"username"`
				Description string `json:"description"`
				Estimation  string `json:"estimation"`
				Date        string `json:"date"`
			}
			var creds CredentialsNewReview

			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err.Error())
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			// Проверяем обязательные поля
			if creds.Username == "" || creds.Description == "" || creds.Estimation == "" || creds.Date == "" {
				log.Printf("Missing required fields")
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			// Создаем отзыв
			est_int, err := strconv.Atoi(creds.Estimation)
			if err != nil {
				log.Printf("Invalid route ID: %v", err)
				http.Error(w, "Invalid route ID", http.StatusBadRequest)
				return
			}
			review := queries.ConstructReview(creds.Username, creds.Description, float64(est_int), creds.Date)

			routeID := r.URL.Query().Get("route_id")
			if routeID == "" {
				// Если route_id отсутствует, перенаправляем на главную страницу
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			routeIDInt, err := strconv.Atoi(routeID)
			if err != nil {
				log.Printf("Invalid route ID: %v", err)
				http.Error(w, "Invalid route ID", http.StatusBadRequest)
				return
			}

			err = review.CreateReview(routeIDInt, data.UserID, p.PoolConns)
			if err != nil {
				log.Printf("Error query CreateNewRoute: %s", err.Error())
				http.Error(w, "Failed to create route", http.StatusInternalServerError)
				return
			}

			// Возвращаем успешный ответ
			http.Redirect(w, r, "/route?route_id="+routeID, http.StatusAccepted)
			return
		}
	}

	return http.HandlerFunc(fn)
}
