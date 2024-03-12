package authHandler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/sabbir-hossain70/api-server/dataHandler"
	"net/http"
	"time"
)

var TokenAuth *jwtauth.JWTAuth

var Secret = []byte("secret_key")

func Login(w http.ResponseWriter, r *http.Request) {
	var cred dataHandler.Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	fmt.Println("Cred ", cred, "   ", err)
	if err != nil {
		http.Error(w, "Can not decode the Credentials... ", http.StatusBadRequest)
		return
	}
	password, ok := dataHandler.CredList[cred.Username]
	if !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if password != cred.Password {
		http.Error(w, "wrong Password", http.StatusBadRequest)
		return
	}
	EndTime := time.Now().Add(10 * time.Minute)

	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{
		"aud": "Sabbir Hossain",
		"exp": EndTime.Unix(),
	})
	fmt.Println("tokenString : ", tokenString)
	if err != nil {
		http.Error(w, "Can not generate jwt", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: EndTime,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully..."))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out"))
}
func InitToken() {
	TokenAuth = jwtauth.New(string(jwa.HS256), Secret, nil)
}
