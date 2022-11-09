package middleware

import (
	"net/http"
	"strings"

	"github.com/emavillamayorpsh/rest-ws/models"
	"github.com/emavillamayorpsh/rest-ws/server"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string	{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

// first param (h) we specify the function that it should go in case everything is okay
// because middleware does a "jump" to the next middleware
func CheckAuthMiddleware(s server.Server) func (h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			// SEND THE PATH TO CHECK IF IT NEEDS TO VALIDATE THE TOKEN
			if !shouldCheckToken(r.URL.Path) {
				// IN CASE EVERYTHING IS CORRECT IT MOVES TO THE NEXT MIDDLEWARE
				next.ServeHTTP(w,r)
				return
			}

			// GET THE TOKEN FROM AUTHORIZATION
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

			// CHECK IF TOKEN IS VALID
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})

			// IN CASE TOKEN INVALID RETURN ERROR
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// IN CASE TOKEN VALID , IT MOVES TO THE NEXT MIDDLEWARE
			next.ServeHTTP(w,r)
		})
	}
}