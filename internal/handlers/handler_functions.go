package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt"
)

type PageData struct {
	Username string
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Claims struct {
	Username string `json:"username"`
	UserID   int    `json: "userID"`
	jwt.StandardClaims
}

func Authorize(a *AppHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Декодируем JSON из тела запроса
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if creds.Email == "" || creds.Password == "" {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Делаем запрос в БД
		conn, err := a.Pool.PoolConns.Acquire(context.Background())
		if err != nil {
			http.Error(w, "Failed to acquire database connection", http.StatusInternalServerError)
			return
		}
		defer conn.Release()

		var userID int
		var user_name string
		err = conn.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1 AND password = $2", creds.Email, creds.Password).Scan(&userID, &user_name)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// ПЕРЕНЕСТИ ВСЕ ЭТО ВОТ ЧТО СНИЗУ В JWT
		// |||||||||||||||||||||||||||||||||||||
		claims := Claims{
			Username: user_name, // Здесь вы можете использовать имя пользователя, полученное из запроса
			UserID:   userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Время истечения токена
			},
		}
		var mySigningKey = []byte("secret")

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		// Создаем куку с JWT
		cookie := http.Cookie{
			Name:     "token",                        // Имя куки
			Value:    tokenString,                    // Значение куки - JWT
			Path:     "/",                            // Путь, для которого кука действительна
			Expires:  time.Now().Add(24 * time.Hour), // Время истечения
			HttpOnly: true,                           // Доступ только через HTTP (не доступен через JavaScript)
			Secure:   false,                          // Установите true, если используете HTTPS
		}

		// Устанавливаем куку в ответ
		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// Обработчик главной страницы
func OpenFirstPage(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/main.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}

func PopularRoutesHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/pages/popular_routes.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}

// Обработчик для страницы "Создать маршрут"
func CreateRouteHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/pages/create_route.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}

// Обработчик для страницы "Мои маршруты"
func ClientRoutesHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/pages/client_route.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}

// Обработчик для страницы "Контакты"
func ContactsHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/pages/contacts.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}

// Обработчик для страницы "О нас"
func AboutUsHandler(a *AppHandlers) http.Handler {
	check_auth := func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			username = ""
		}

		// Если пользователь авторизирован - ему выставляется имя
		tmpl, err := template.ParseFiles("./web/pages/about_us.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Username: username,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(check_auth)
}
