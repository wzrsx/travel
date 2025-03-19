package handler_functions

import (
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
	RouteID           int
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
		// Обработка GET-запроса
		if r.Method == http.MethodGet {
			routeID := r.URL.Query().Get("route_id")
			if routeID == "" {
				http.Error(w, "Route ID is required", http.StatusBadRequest)
				return
			}
			
			routeIDInt, err := strconv.Atoi(routeID)
			if err != nil {
				log.Printf("Invalid route ID: %v", err)
				http.Error(w, "Invalid route ID", http.StatusBadRequest)
				return
			}

			// Получаем информацию о маршруте
			queryTakeRoute, err := queries.TakeRouteInfo(routeIDInt, p.PoolConns)
			if err != nil {
				log.Printf("Error Take Route info: %v", err)
				http.Error(w, "Error Take Route info", http.StatusInternalServerError)
				return
			}

			// Проверяем, есть ли уже параметры маршрута в URL
			if r.URL.Query().Get("route_name") == "" || r.URL.Query().Get("route_place") == "" || r.URL.Query().Get("routeLink") == "" || r.URL.Query().Get("route_description") == "" || r.URL.Query().Get("route_estimation") == "" {
				// Формируем URL с параметрами
				redirectURL := fmt.Sprintf("/route?route_id=%s&route_name=%s&route_place=%s&routeLink=%s&route_description=%s&route_estimation=%s",
					url.QueryEscape(routeID),
					url.QueryEscape(queryTakeRoute.Route_name),
					url.QueryEscape(queryTakeRoute.Route_place),
					url.QueryEscape(queryTakeRoute.Yandex_route),
					url.QueryEscape(queryTakeRoute.Route_description),
					url.QueryEscape(strconv.Itoa(queryTakeRoute.Route_estimation)),
				)

				// Перенаправляем с параметрами
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			// Если параметры уже есть, продолжаем обработку GET-запроса
			tmpl := template.Must(template.New("route.html").Funcs(funcMap).ParseFiles("web/pages/route.html"))
			data, err := check_username_data(w, r)
			if err != nil {
				log.Printf("Failed to check username: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Получение данных маршрута
			photos, err := queries.CreatePhotos(routeIDInt, p.PoolConns)
			if err != nil {
				log.Printf("Error getting photos: %s", err.Error())
				http.Error(w, "Error getting photos", http.StatusInternalServerError)
				return
			}

			reviews, err := queries.TakeReviews(routeIDInt, p.PoolConns)
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

			// Получение параметров маршрута из URL
			dataRoute := DataRoute{
				UserID:            data.UserID,
				Username:          data.Username,
				RouteID:           routeIDInt,
				RouteLink:         decodedLink,
				Route_name:        r.URL.Query().Get("route_name"),
				Route_place:       r.URL.Query().Get("route_place"),
				Route_description: r.URL.Query().Get("route_description"),
				Route_estimation:  r.URL.Query().Get("route_estimation"),
				Route_photos:      photos,
				Route_reviews:     reviews,
			}

			err = tmpl.Execute(w, dataRoute)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// Обработка POST-запроса
		if r.Method == http.MethodPost {
			routeID := r.URL.Query().Get("route_id")
			routeIDInt, err := strconv.Atoi(routeID)
			if err != nil {
				log.Printf("Invalid route ID: %v", err)
				http.Error(w, "Invalid route ID", http.StatusBadRequest)
				return
			}

			queryTakeRoute, err := queries.TakeRouteInfo(routeIDInt, p.PoolConns)
			if err != nil {
				log.Printf("Error Take Route info: %v", err)
				http.Error(w, "Error Take Route info:", http.StatusBadRequest)
				return
			}

			// Формируем URL с параметрами
			redirectURL := fmt.Sprintf("/route?route_id=%s&route_name=%s&route_place=%s&routeLink=%s&route_description=%s&route_estimation=%s",
				url.QueryEscape(routeID),
				url.QueryEscape(queryTakeRoute.Route_name),
				url.QueryEscape(queryTakeRoute.Route_place),
				url.QueryEscape(queryTakeRoute.Yandex_route),
				url.QueryEscape(queryTakeRoute.Route_description),
				url.QueryEscape(strconv.Itoa(queryTakeRoute.Route_estimation)),
			)

			// Перенаправляем с параметрами
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
	}

	return http.HandlerFunc(check_auth)
}
