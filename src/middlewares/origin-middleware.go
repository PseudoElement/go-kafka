package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pseudoelement/go-kafka/src/shared"
)

func OriginMiddleware(next http.Handler) http.Handler {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(originsEnv, " ")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Println("allowedOrigins ==>", origins)
		log.Println("oridgin ==>", origin)
		if !shared.Contains(origins, origin) {
			msg := "Origin " + origin + " is not allowed."
			shared.FailedResp(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
