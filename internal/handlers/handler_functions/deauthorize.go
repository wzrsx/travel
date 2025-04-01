package handler_functions

import "net/http"

func DeAuthorize() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type Cookies struct {
			JWT *http.Cookie
		}
		var cookies Cookies
		cookies.JWT = &http.Cookie{
			Name:     "token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
		}
		http.SetCookie(w, cookies.JWT)
	}
	return http.HandlerFunc(fn)
}
