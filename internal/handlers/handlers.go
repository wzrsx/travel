package handlers

import (
	"net/http"
	"travel/internal/handlers/JWT"
	"travel/internal/handlers/handler_functions"
	"travel/internal/repositories/pool_conections"
)

type codeData struct {
	Email string
	Code  string
}

type AppHandlers struct {
	Pool *pool_conections.Pool_conections
}

func NewAppHandlers(pool *pool_conections.Pool_conections) *AppHandlers {
	return &AppHandlers{Pool: pool}
}

func (a *AppHandlers) SetHandlers() {
	http.Handle("/", JWT.JWTMiddleware(handler_functions.OpenFirstPage()))
	http.Handle("/popular-routes", JWT.JWTMiddleware(handler_functions.PopularRoutesHandler(a.Pool)))
	http.Handle("/create-route", JWT.JWTMiddleware(handler_functions.CreateRouteHandler(a.Pool)))
	http.Handle("/client-routes", JWT.JWTMiddleware(handler_functions.ClientRoutesHandler(a.Pool)))
	http.Handle("/contacts", JWT.JWTMiddleware(handler_functions.ContactsHandler(a.Pool)))
	http.Handle("/about-us", JWT.JWTMiddleware(handler_functions.AboutUsHandler(a.Pool)))
	http.Handle("/route", JWT.JWTMiddleware(handler_functions.OpenRoutePage(a.Pool)))
	http.Handle("/route/reviews", JWT.JWTMiddleware(handler_functions.ReviewsHandler(a.Pool)))
	http.Handle("/route/review", JWT.JWTMiddleware(handler_functions.ReviewHandler(a.Pool)))



	http.Handle("/save-route", JWT.JWTMiddleware(handler_functions.CreateRouteHandler(a.Pool)))

	http.Handle("/authorize", handler_functions.Authorize(a.Pool))
	http.Handle("/registration", handler_functions.Registration(a.Pool))
	http.Handle("/deauthorize", handler_functions.DeAuthorize())

	http.Handle("/send_to_email/pass_code_recovery", handler_functions.SendEmailMessageWithCodeRecovery())
	http.Handle("/send_to_email/pass_code_registration", handler_functions.SendEmailMessageWithCodeRecovery())

	http.Handle("/check/email", handler_functions.CheckEmailIntoDB(a.Pool))
	http.Handle("/check/pass_code", handler_functions.CheckPassCode())

	http.Handle("/change/password", JWT.JWTMiddleware(handler_functions.ChangePassword(a.Pool)))

}

func (a *AppHandlers) SetDirs() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./web/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./web/scripts"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./web/pages"))))
	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("./web/photos"))))
}
