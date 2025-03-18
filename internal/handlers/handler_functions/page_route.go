package handler_functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

type DataRoute struct {
	UserID            int
	Username          string
	RouteLink         string
	Route_name        string
	Route_place       string
	Route_description string
	Route_estimation  string
	Route_photos      []queries.Photo
	Route_reviews     []queries.Review
}

func OpenRoutePage(p *pool_conections.Pool_conections) http.Handler {
	funcMap := template.FuncMap{
		"seq": func(n int) []int {
			var sequence []int
			for i := 1; i <= n; i++ {
				sequence = append(sequence, i)
			}
			return sequence
		},
	}
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Обработка POST-запроса
			var requestData struct {
				RouteName        string `json:"route_name"`
				RoutePlace       string `json:"route_place"`
				RouteLink        string `json:"routeLink"`
				RouteDescription string `json:"route_description"`
				RouteEstimation  string `json:"route_estimation"`
			}

			if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
				log.Printf("Failed to decode JSON: %v", err)
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			routeID := r.URL.Query().Get("route_id")

			// Формируем URL с параметрами
			redirectURL := fmt.Sprintf("/route?route_id=%s&route_name=%s&route_place=%s&routeLink=%s&route_description=%s&route_estimation=%s",
				url.QueryEscape(routeID),
				url.QueryEscape(requestData.RouteName),
				url.QueryEscape(requestData.RoutePlace),
				url.QueryEscape(requestData.RouteLink),
				url.QueryEscape(requestData.RouteDescription),
				url.QueryEscape(requestData.RouteEstimation),
			)

			// Перенаправляем с параметрами
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		// Обработка GET-запроса (отображение страницы)
		tmpl := template.Must(template.New("route.html").Funcs(funcMap).ParseFiles("web/pages/route.html"))
		data, err := check_username_data(w, r)
		if err != nil {
			log.Printf("Failed to check username: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		routeID := r.URL.Query().Get("route_id")
		routeIDInt, err := strconv.Atoi(routeID)
		if err != nil {
			log.Printf("Invalid route ID: %v", err)
			http.Error(w, "Invalid route ID", http.StatusBadRequest)
			return
		}

		// Получение данных маршрута
		// Фото
		photos, err := queries.CreatePhotos(routeIDInt, p.PoolConns)
		if err != nil {
			log.Printf("Error getting photos: %s", err.Error())
			http.Error(w, "Error getting photos", http.StatusInternalServerError)
			return
		}

		// Отзывы
		reviews, err := queries.CreateReviews(routeIDInt, p.PoolConns)
		if err != nil {
			log.Printf("Error getting reviews: %s", err.Error())
			http.Error(w, "Error getting reviews", http.StatusInternalServerError)
			return
		}

		// Декодирование ссылки на карты
		encodedLink := r.URL.Query().Get("routeLink")
		decodedLink, err := url.QueryUnescape(encodedLink) // Первое декодирование
		if err != nil {
			log.Printf("Failed to decode map link (first pass): %s", err.Error())
			http.Error(w, "Failed to decode map link", http.StatusInternalServerError)
			return
		}

		decodedLink, err = url.QueryUnescape(decodedLink) // Второе декодирование
		if err != nil {
			log.Printf("Failed to decode map link (second pass): %s", err.Error())
			http.Error(w, "Failed to decode map link", http.StatusInternalServerError)
			return
		}

		log.Println(decodedLink)

		// Получение параметров маршрута из URL
		dataRoute := DataRoute{
			UserID:            data.UserID,
			Username:          data.Username,
			RouteLink:         decodedLink,
			Route_name:        r.URL.Query().Get("route_name"),
			Route_place:       r.URL.Query().Get("route_place"),
			Route_description: r.URL.Query().Get("route_description"),
			Route_estimation:  r.URL.Query().Get("route_estimation"),
			Route_photos:      photos,
			Route_reviews:     reviews,
		}
		log.Print(dataRoute.Route_photos)

		err = tmpl.Execute(w, dataRoute)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}
