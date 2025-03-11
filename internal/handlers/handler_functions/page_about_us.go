package handler_functions

import (
	"log"
	"net/http"
	"travel/internal/repositories/pool_conections"
)

// Обработчик для страницы "О нас"
func AboutUsHandler(p *pool_conections.Pool_conections) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/about_us.html")
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Failed to render template: %s", err.Error())
			return
		}
	}
	return http.HandlerFunc(check_auth)
}
