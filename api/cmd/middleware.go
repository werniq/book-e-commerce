package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func (app *application) IsAuthorized(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			var err error
			err = errors.New("no token found")
			json.NewEncoder(w).Encode(err)
		}

		var mySigningKey = os.Getenv("JWT_AUTH_SECRET_KEY")
		// jwt.Parse(tokenString string, keyFunc func(*Token) Name() e
		// Parse, validate, and return a token.
		// keyFunc will receive the parsed token and should return the key for validating.
		// If everything is kosher, err will be nil
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			// Parse methods use this callback function to supply
			// the key for verification.  The function receives the parsed,
			// but unverified Token.  This allows you to use properties in the
			// Header of the token (such as `kid`) to identify which key to use.
			//type Keyfunc func(*Token) (interface{}, error)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("there was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(errors.New("your token has been expired"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("role", "admin")
				handler.ServeHTTP(w, r)
				json.NewEncoder(w).Encode(token)
				//app.database.GetUserForToken()
				return
			} else if claims["role"] == "user" {
				r.Header.Set("role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}
		json.NewEncoder(w).Encode(errors.New("unauthorized"))
		//http.Redirect(w, r, "localhost:3000/register", http.StatusUnauthorized)
	})

}
