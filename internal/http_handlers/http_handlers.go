package http_handlers

import "net/http"

// Обработчик для главной страницы
func OpenFirstPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/main.html")
}

func SetDirs() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./web/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./web/scripts"))))
}
