package login

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"huru/models/person"
	"huru/mw"
	"log"
	"net/http"
	"time"
)

// Handler just what it says
type Handler struct{}

type JwtToken struct {
	Token string `json:"token"`
}

type Injection struct {
	PeopleById, PeopleByHandle, PeopleByEmail person.Map
}

// Mount just what it says
func (h Handler) Mount(router *mux.Router, v Injection) Handler {

	mwList := mw.List(mw.Error())
	router.HandleFunc("/login", mw.Middleware(mwList, h.makeLogin(v))).Methods("POST")
	return h
}

// MakeGetMany Display all from the people var
func (h Handler) makeLogin(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var user person.Model;
		json.NewDecoder(r.Body).Decode(&user)

		fmt.Printf("%+v\n", user)

		//if _, ok := v.PeopleByHandle[user.Handle]; !ok {
		//	panic("Please register b/c you cannot login because no user exists with handle " + user.Handle)
		//}

		if _, ok := v.PeopleByEmail[user.Email]; !ok {
			panic("Please register b/c you cannot login because no user exists with this email" + user.Email)
			return;
		}


		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			//"username": user.Handle,
			//"handle":   user.Handle,
			"password": user.Password,
			"email":    user.Email,
			"role":     user.Role,
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {
			panic(err);
		}

		log.Println("token string", tokenString)

		expire := time.Now().Add(5000 * time.Minute)
		cookie := http.Cookie{Name: "x-huru-api-token", Value: tokenString, Path: "/", Expires: expire, MaxAge: 90000}

		log.Println("Creating new cookie.")

		http.SetCookie(w, &cookie)

		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	})
}
