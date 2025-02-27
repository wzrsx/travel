package http_handlers

import "net/http"

// Обработчик для главной страницы
func OpenFirstPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/main.html")
}

// Обработчик для страницы "Популярные маршруты"
func PopularRoutesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/pages/popular_routes.html")
}

// Обработчик для страницы "Создать маршрут"
func CreateRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/pages/create_route.html")
}

// Обработчик для страницы "Мои маршруты"
func ClientRoutesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/pages/client_routes.html")
}

// Обработчик для страницы "Контакты"
func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/pages/contacts.html")
}

// Обработчик для страницы "О нас"
func AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/pages/about_us.html")
}

