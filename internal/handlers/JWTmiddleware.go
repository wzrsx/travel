package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из cookies
		cookie, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Парсим токен
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Fprintf(w, "unexpected signing method: %v", token.Header["alg"])
			return
		}
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			// http.Redirect(w, r, "/authorize", http.StatusSeeOther)
			return
		}

		// Сохраняем данные пользователя в контекст запроса
		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "username", claims["username"])
		ctx = context.WithValue(ctx, "email", claims["email"])
		r = r.WithContext(ctx)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
