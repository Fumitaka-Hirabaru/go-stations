package middleware

import (
	"net/http"
	"os"
)

func getCredentials() (string, string) {
	userID := os.Getenv("BASIC_AUTH_USER_ID")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	return userID, password
}

func BasicAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, password := getCredentials()
		user, pass, ok := r.BasicAuth()

		if !ok || user != userID || pass != password {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
