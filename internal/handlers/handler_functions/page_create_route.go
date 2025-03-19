package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

// Обработчик для страницы "Создать маршрут"
func CreateRouteHandler(p *pool_conections.Pool_conections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/create_route.html")
		if err != nil {
			log.Printf("%s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if data.UserID == 0 {
			http.Redirect(w, r, "/?openLoginDialog=true", http.StatusSeeOther)
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
			type CredentialsCreateRoute struct {
				Yandex_route      string `json:"routeLink"`
				Route_name        string `json:"route_name"`
				Route_place       string `json:"route_place"`
				Route_description string `json:"route_description"`
			}
			var creds CredentialsCreateRoute

			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err.Error())
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			// Проверяем обязательные поля
			if creds.Yandex_route == "" || creds.Route_name == "" || creds.Route_place == "" || creds.Route_description == "" {
				log.Printf("Missing required fields")
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			// Создаем маршрут
			route := queries.ConstructRoute(creds.Yandex_route, creds.Route_name, creds.Route_place, creds.Route_description, data.UserID)
			err = route.CreateNewRoute(p.PoolConns)
			if err != nil {
				log.Printf("Error query CreateNewRoute: %s", err.Error())
				http.Error(w, "Failed to create route", http.StatusInternalServerError)
				return
			}

			// Возвращаем успешный ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Route created successfully"})
			return
		}

		// Обработка GET-запроса
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Failed to render template: %s", err.Error())
		}
	}
	return http.HandlerFunc(fn)
}
