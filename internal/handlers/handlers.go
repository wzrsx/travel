package handlers

import (
	"net/http"
	"travel/internal/http_handlers"
)

func SetHandlers(){
	http.HandleFunc("/", http_handlers.OpenFirstPage)
	http.HandleFunc("/popular-routes", http_handlers.PopularRoutesHandler)
	http.HandleFunc("/create-route", http_handlers.CreateRouteHandler)
	http.HandleFunc("/client-routes", http_handlers.ClientRoutesHandler)
	http.HandleFunc("/contacts", http_handlers.ContactsHandler)
	http.HandleFunc("/about-us", http_handlers.AboutUsHandler)
}

func SetDirs() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./web/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./web/scripts"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./web/pages"))))
}

