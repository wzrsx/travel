package handlers

import (
	"net/http"
	"travel/internal/handlers/JWT"
	"travel/internal/repositories/pool_conections"
)

type AppHandlers struct {
	Pool *pool_conections.Pool_conections
}

func NewAppHandlers(pool *pool_conections.Pool_conections) *AppHandlers {
	return &AppHandlers{Pool: pool}
}

func (a *AppHandlers) SetHandlers() {
	http.Handle("/", JWT.JWTMiddleware(OpenFirstPage(a)))
	http.Handle("/popular-routes", JWT.JWTMiddleware(PopularRoutesHandler(a)))
	http.Handle("/create-route", JWT.JWTMiddleware(CreateRouteHandler(a)))
	http.Handle("/client-routes", JWT.JWTMiddleware(ClientRoutesHandler(a)))
	http.Handle("/contacts", JWT.JWTMiddleware(ContactsHandler(a)))
	http.Handle("/about-us", JWT.JWTMiddleware(AboutUsHandler(a)))
	http.Handle("/authorize", Authorize(a))
	http.Handle("/registration", Registration(a))

}

func (a *AppHandlers) SetDirs() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./web/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./web/scripts"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./web/pages"))))
}
