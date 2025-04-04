package handler_functions

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

func ReviewsHandler(p *pool_conections.Pool_conections) http.Handler {
	funcMap := template.FuncMap{
		"seq": func(n float64) []float64 {
			var sequence []float64
			for i := 1.00; i <= n; i++ {
				sequence = append(sequence, i)
			}
			return sequence
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.New("reviews.html").Funcs(funcMap).ParseFiles("web/pages/reviews.html"))

			// Получаем данные пользователя
			data, err := check_username_data(r)
			if err != nil {
				log.Printf("Failed to check username: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Проверяем наличие route_id в URL
			routeID := r.URL.Query().Get("route_id")
			if routeID == "" {
				// Если route_id отсутствует, перенаправляем на главную страницу
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			routeIdINT, err := strconv.Atoi(routeID)
			if err != nil {
				http.Error(w, "Invalid route_id format", http.StatusBadRequest)
				return
			} 

			reviews, err := queries.TakeReviews(routeIdINT, p.PoolConns)
			if err != nil{
				log.Printf("Error query review: %v", err)
				return
			}
			
			// Подготавливаем данные для шаблона
			dataWithReviews := struct {
				Username string
				UserID   int
				Reviews  []queries.Review
			}{
				Username: data.Username,
				UserID:   data.UserID,
				Reviews:  reviews,
			}

			// Рендерим шаблон
			err = tmpl.Execute(w, dataWithReviews)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				log.Printf("Failed to render template: %s", err.Error())
			}
			return
		}
	}

	return http.HandlerFunc(fn)
}
