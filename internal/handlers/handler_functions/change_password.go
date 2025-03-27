package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
)

func ChangePassword() http.Handler{
	fn := func (w http.ResponseWriter, r *http.Request){
		if(r.Method == http.MethodPost){
			type Credentials struct{
				Email string `json:"email"`
			}
			var creds Credentials
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err!=nil{
				log.Printf("Error decode JSON: %v", err)
				return
			}

			if(!canChange[creds.Email]){
				w.WriteHeader(http.StatusConflict)
			}
		}
	}
	return http.HandlerFunc(fn)
}