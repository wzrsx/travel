package JWT

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"userID"`
	jwt.StandardClaims
}

func SetCookeJWT(w http.ResponseWriter, userID int, username string) error {
	claims := Claims{
		Username: username,
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
		return err
	}

	// Создаем куку с JWT
	cookieJWT := http.Cookie{
		Name:     "token",                        // Имя куки
		Value:    tokenString,                    // Значение куки - JWT
		Path:     "/",                            // Путь, для которого кука действительна
		Expires:  time.Now().Add(24 * time.Hour), // Время истечения
		HttpOnly: true,                           // Доступ только через HTTP (не доступен через JavaScript)
		Secure:   false,                          // Установите true, если используете HTTPS
	}

	// Устанавливаем куку в ответ
	http.SetCookie(w, &cookieJWT)
	return nil
}
