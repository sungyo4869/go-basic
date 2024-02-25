package middleware

import (
	"net/http"
	"os"
	"log"
	"github.com/joho/godotenv"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		err := godotenv.Load()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	  
		userID := os.Getenv("BASIC_AUTH_USER_ID")
		password := os.Getenv("BASIC_AUTH_USER_PASSWORD")

		log.Println(userID)
		log.Println(password)

		ui, pass, ok := r.BasicAuth()
 
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if userID != ui || password != pass {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
