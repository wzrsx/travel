package handlers

import (
	"net/http"
	"travel/internal/repositories/pool_conections"
)

type AppHandlers struct {
	Pool *pool_conections.Pool_conections
}

func NewAppHandlers(pool *pool_conections.Pool_conections) *AppHandlers {
	return &AppHandlers{Pool: pool}
}

func (a *AppHandlers) SetHandlers() {
	http.Handle("/", JWTMiddleware(OpenFirstPage(a)))
	http.Handle("/popular-routes", JWTMiddleware(PopularRoutesHandler(a)))
	http.Handle("/create-route", JWTMiddleware(CreateRouteHandler(a)))
	http.Handle("/client-routes", JWTMiddleware(ClientRoutesHandler(a)))
	http.Handle("/contacts", JWTMiddleware(ContactsHandler(a)))
	http.Handle("/about-us", JWTMiddleware(AboutUsHandler(a)))
	http.Handle("/authorize", Authorize(a))
	
}

func (a *AppHandlers) SetDirs() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./web/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./web/scripts"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./web/pages"))))
}
