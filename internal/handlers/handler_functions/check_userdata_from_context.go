package handler_functions

import (
	"fmt"
	"net/http"
	"text/template"
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

func check_username_data(w http.ResponseWriter, r *http.Request) (*PageData, error) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		username = ""
	}
	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		user_id = 0
	}

	data := PageData{
		Username: username,
		UserID:   user_id,
	}

	return &data, nil

}