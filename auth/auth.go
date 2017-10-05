package auth

import (
	"fmt"
	"log"
	"net/http"

	m "github.com/adam72m/go-web/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("dwadziescia-muharadzinow-bije-trzech-rabinow"))

func SessionCheckingHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, fmt.Sprintf("The session expired \n"), http.StatusInternalServerError)
			return
		}
		userName := session.Values["user"]
		log.Printf("%v", userName)
		if userName != "adam" {
			http.Error(w, fmt.Sprintf("Invalid login \n"), http.StatusInternalServerError)
			return
		}

		fn(w, r)
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no Auth cookie is set then return a 404 not found
		cookie, err := r.Cookie("Auth")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Return a Token using the cookie
		token, err := jwt.ParseWithClaims(cookie.Value, &m.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
