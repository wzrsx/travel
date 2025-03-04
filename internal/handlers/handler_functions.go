package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"travel/internal/handlers/JWT"
	"travel/internal/repositories/queries"
)

type PageData struct {
	Username string
	UserID   int
}

func check_username(w http.ResponseWriter, r *http.Request, web_file string) (*template.Template, *PageData, error) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		username = ""
	}
	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		user_id = 0
	}

	// Если пользователь авторизирован - ему выставляется имя
	tmpl, err := template.ParseFiles(web_file)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return nil, nil, fmt.Errorf("Failed to load template: %s", err.Error())
	}

	data := PageData{
		Username: username,
		UserID:   user_id,
	}

	return tmpl, &data, nil

}

func Authorize(a *AppHandlers) http.Handler {
	authorize := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем Content-Type
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "Invalid Content-Type: ", http.StatusUnsupportedMediaType)
			return
		}

		// Декодируем JSON из тела запроса
		type CredentialsAuthorize struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var creds CredentialsAuthorize
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if creds.Email == "" || creds.Password == "" {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Делаем запрос в БД
		user_auth_res := queries.NewUserAuthResult()
		err = user_auth_res.AuthorizeQuery(creds.Email, creds.Password, a.Pool.PoolConns)
		if err != nil {
			log.Printf("Authorize Querry returns error: %v", err)
			http.Error(w, "Authorize Querry error", http.StatusBadRequest)
		}

		// устанавливаем Cookie с JWT
		err = JWT.SetCookeJWT(w, user_auth_res.UserID, user_auth_res.Username)
		if err != nil {
			http.Error(w, "Error set cookie JWT", http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(authorize)
}

func Registration(a *AppHandlers) http.Handler {
	registration := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "Invalid Content-Type: ", http.StatusUnsupportedMediaType)
			return
		}

		// Декодируем JSON из тела запроса
		type CredentialsRegistration struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var creds CredentialsRegistration
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if creds.Email == "" || creds.Password == "" || creds.Username == "" {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Делаем запрос в БД
		user_reg_res := queries.NewUserRegistrationResult()
		err = user_reg_res.RegistrationQuery(creds.Username, creds.Email, creds.Password, a.Pool.PoolConns)
		if err != nil {
			log.Printf("Registration Querry returns error: %v", err)
			http.Error(w, "Registration Querry error", http.StatusBadRequest)
		}

		// устанавливаем Cookie с JWT
		err = JWT.SetCookeJWT(w, user_reg_res.UserID, user_reg_res.Username)
		if err != nil {
			http.Error(w, "Error set cookie JWT", http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(registration)
}

// Обработчик главной страницы
func OpenFirstPage(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/main.html")
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

func PopularRoutesHandler(a *AppHandlers) http.Handler {
	routes, err := queries.TakeRoutesInfoQueryPopular(a.Pool.PoolConns)
	if err != nil {
		log.Printf("Error routes having: %v", err)
	}
	for _, route := range routes {
		log.Printf("Маршрут ID: %d", route.IdRoute)
		log.Printf("Название: %s", route.Name_route)
		log.Printf("Ссылка на Yandex: %s", route.Yandex_route)
		log.Printf("Оценка: %d", route.Estimation)
		log.Printf("Фото: %s", route.PathToPhoto)

		// Вывод отзывов
		log.Println("Отзывы:")
		for _, review := range route.Reviews {
			log.Printf("  Описание: %s", review.Description)
			log.Printf("  Оценка: %d", review.Estimation)
		}
		log.Println("-----------------------------")
	}

	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/popular_routes.html")
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

// Обработчик для страницы "Создать маршрут"
func CreateRouteHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		check_username(w, r, "web/pages/create_route.html")
		tmpl, data, err := check_username(w, r, "web/pages/create_route.html")
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		if r.Method == http.MethodPost {
			contentType := r.Header.Get("Content-type")
			if contentType != "application/json" {
				http.Error(w, "Invalid Content-Type: ", http.StatusUnsupportedMediaType)
				return
			}

			type CredentialsCreateRoute struct {
				Yandex_route      string `json: "routeLink"`
				Route_name        string `json: "route_name"`
				Route_place       string `json: "route_place"`
				Route_description string `json: "route_description"`
			}
			var creds CredentialsCreateRoute

			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				log.Printf("Error decoding json by creating route: %s", err.Error())
				return
			}
			if creds.Yandex_route == "" || creds.Route_name == "" || creds.Route_place == "" || creds.Route_description == "" {
				log.Printf("Error decoding JSON: %s", err.Error())
				http.Error(w, "Yandex_route is required", http.StatusBadRequest)
				return
			}

			route := queries.CreateRouteStruct(creds.Yandex_route, creds.Route_name, creds.Route_place, creds.Route_description, data.UserID)
			err = route.CreateNewRoute(a.Pool.PoolConns)
			if err != nil {
				log.Printf("Error query CreateNewRoute: %s", err.Error())
				http.Error(w, "Error query CreateNewRoute", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
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

// Обработчик для страницы "Мои маршруты"
func ClientRoutesHandler(a *AppHandlers) http.Handler {

	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/client_routes.html")
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}

		// Logic with take UsersRoutes
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Println("error not ok")
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
			return
		}

		routes, err := queries.TakeUsersRoutesInfoQuery(a.Pool.PoolConns, userID)
		if err != nil {
			log.Printf("Error routes having: %v", err)
		}

		for _, route := range routes {
			log.Printf("Маршрут ID: %d", route.IdRoute)
			log.Printf("Название: %s", route.Name_route)
			log.Printf("Ссылка на Yandex: %s", route.Yandex_route)
			log.Printf("Оценка: %d", route.Estimation)
			log.Printf("Фото: %s", route.PathToPhoto)

			// Вывод отзывов
			log.Println("Отзывы:")
			for _, review := range route.Reviews {
				log.Printf("  Описание: %s", review.Description)
				log.Printf("  Оценка: %d", review.Estimation)
			}
			log.Println("-----------------------------")
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

// Обработчик для страницы "Контакты"
func ContactsHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/contacts.html")
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

// Обработчик для страницы "О нас"
func AboutUsHandler(a *AppHandlers) http.Handler {
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
