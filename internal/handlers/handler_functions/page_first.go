package handler_functions

import (
	"log"
	"net/http"
	"travel/internal/repositories/pool_conections"
)

// Обработчик главной страницы
func OpenFirstPage(a *pool_conections.Pool_conections) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/main.html")
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		type PgData struct {
			UserID    int
			Username  string
			IsNotAuth bool
		}
		pg_data := PgData{
			UserID:    data.UserID,
			Username:  data.Username,
			IsNotAuth: data.UserID == 0,
		}
		err = tmpl.Execute(w, pg_data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Failed to render template: %s", err.Error())
			return
		}
	}
	return http.HandlerFunc(check_auth)
}
