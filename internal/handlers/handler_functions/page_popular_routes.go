package handler_functions

import (
	"log"
	"net/http"
	"text/template"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func PopularRoutesHandler(p *pool_conections.Pool_conections) http.Handler {
	// Определяем функцию seq
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
		// Парсим шаблон с функцией seq
		tmpl := template.Must(template.New("popular_routes.html").Funcs(funcMap).ParseFiles("web/pages/popular_routes.html"))

		// Получаем данные пользователя
		data, err := check_username_data(r)
		if err != nil {
			log.Printf("Failed to check username: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Получаем данные маршрутов
		routes, err := queries.TakeRoutesInfoQueryPopular(p.PoolConns)
		if err != nil {
			log.Printf("Failed to get routes data: %v", err)
			http.Error(w, "Failed to fetch routes", http.StatusInternalServerError)
			return
		}

		// Подготавливаем данные для шаблона
		dataWithRoutes := struct {
			Username   string
			UserID     int
			RoutesInfo []queries.RoutesInfoResults
		}{
			Username:   data.Username,
			UserID:     data.UserID,
			RoutesInfo: routes,
		}

		// Рендерим шаблон
		err = tmpl.Execute(w, dataWithRoutes)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(check_auth)
}